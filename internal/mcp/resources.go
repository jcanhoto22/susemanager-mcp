package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	mcp_golang "github.com/metoro-io/mcp-golang"
)

func (s *Server) registerResources() error {
	resources := []struct {
		uri         string
		name        string
		description string
		mimeType    string
		handler     func() (*mcp_golang.ResourceResponse, error)
	}{
		{
			uri:         "suma://systems",
			name:        "All Systems",
			description: "List of all registered systems in SUSE Manager",
			mimeType:    "application/json",
			handler:     s.handleSystemsResource(),
		},
		{
			uri:         "suma://channels",
			name:        "All Channels",
			description: "List of all software channels in SUSE Manager",
			mimeType:    "application/json",
			handler:     s.handleChannelsResource(),
		},
		{
			uri:         "suma://summary",
			name:        "SUSE Manager Summary",
			description: "Summary overview of SUSE Manager (system count, active systems, etc.)",
			mimeType:    "application/json",
			handler:     s.handleSummaryResource(),
		},
	}

	for _, r := range resources {
		if err := s.mcp.RegisterResource(r.uri, r.name, r.description, r.mimeType, r.handler); err != nil {
			return fmt.Errorf("register resource %s: %w", r.uri, err)
		}
		s.logger.Debug("resource registered", "uri", r.uri)
	}

	s.logger.Info("resources registered", "count", len(resources))
	return nil
}

func (s *Server) handleSystemsResource() func() (*mcp_golang.ResourceResponse, error) {
	svc := NewServices(s.api)
	return func() (*mcp_golang.ResourceResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		systems, err := svc.System.ListSystems(ctx)
		if err != nil {
			return mcp_golang.NewResourceResponse(
				mcp_golang.NewTextEmbeddedResource(
					"suma://systems",
					fmt.Sprintf(`{"error":"%s"}`, err.Error()),
					"application/json",
				),
			), nil
		}

		data, _ := json.MarshalIndent(systems, "", "  ")
		return mcp_golang.NewResourceResponse(
			mcp_golang.NewTextEmbeddedResource(
				"suma://systems",
				string(data),
				"application/json",
			),
		), nil
	}
}

func (s *Server) handleChannelsResource() func() (*mcp_golang.ResourceResponse, error) {
	svc := NewServices(s.api)
	return func() (*mcp_golang.ResourceResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		channels, err := svc.Channel.ListAllChannels(ctx)
		if err != nil {
			return mcp_golang.NewResourceResponse(
				mcp_golang.NewTextEmbeddedResource(
					"suma://channels",
					fmt.Sprintf(`{"error":"%s"}`, err.Error()),
					"application/json",
				),
			), nil
		}

		data, _ := json.MarshalIndent(channels, "", "  ")
		return mcp_golang.NewResourceResponse(
			mcp_golang.NewTextEmbeddedResource(
				"suma://channels",
				string(data),
				"application/json",
			),
		), nil
	}
}

func (s *Server) handleSummaryResource() func() (*mcp_golang.ResourceResponse, error) {
	svc := NewServices(s.api)
	return func() (*mcp_golang.ResourceResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		systems, sysErr := svc.System.ListSystems(ctx)
		channels, chErr := svc.Channel.ListAllChannels(ctx)
		active, actErr := svc.System.ListActiveSystems(ctx)

		summary := map[string]interface{}{
			"total_systems":     len(systems),
			"active_systems":    len(active),
			"total_channels":    len(channels),
			"errors":            map[string]string{},
		}

		if sysErr != nil {
			summary["errors"].(map[string]string)["systems"] = sysErr.Error()
		}
		if chErr != nil {
			summary["errors"].(map[string]string)["channels"] = chErr.Error()
		}
		if actErr != nil {
			summary["errors"].(map[string]string)["active"] = actErr.Error()
		}

		data, _ := json.MarshalIndent(summary, "", "  ")
		return mcp_golang.NewResourceResponse(
			mcp_golang.NewTextEmbeddedResource(
				"suma://summary",
				string(data),
				"application/json",
			),
		), nil
	}
}
