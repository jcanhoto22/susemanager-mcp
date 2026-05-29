package mcp

import (
	"context"
	"fmt"
	"log/slog"

	mcp_golang "github.com/metoro-io/mcp-golang"
	mcphttp "github.com/metoro-io/mcp-golang/transport/http"
	"github.com/metoro-io/mcp-golang/transport/stdio"

	"mcp-susemanager/internal/api"
	"mcp-susemanager/internal/config"
)

type Server struct {
	mcp    *mcp_golang.Server
	logger *slog.Logger
	cfg    *config.Config
	api    *api.Client
}

func New(cfg *config.Config, apiClient *api.Client, logger *slog.Logger) (*Server, error) {
	s := &Server{
		logger: logger.With("component", "mcp-server"),
		cfg:    cfg,
		api:    apiClient,
	}

	if err := s.setupTransport(); err != nil {
		return nil, fmt.Errorf("setup transport: %w", err)
	}

	if err := s.registerTools(); err != nil {
		return nil, fmt.Errorf("register tools: %w", err)
	}

	if err := s.registerResources(); err != nil {
		return nil, fmt.Errorf("register resources: %w", err)
	}

	if err := s.registerPrompts(); err != nil {
		return nil, fmt.Errorf("register prompts: %w", err)
	}

	return s, nil
}

func (s *Server) setupTransport() error {
	switch s.cfg.Server.Transport {
	case "http":
		transport := mcphttp.NewHTTPTransport(s.cfg.Server.MCPEndpoint)
		transport.WithAddr(fmt.Sprintf(":%d", s.cfg.Server.Port))
		s.mcp = mcp_golang.NewServer(transport)
		s.logger.Info("MCP HTTP transport configured",
			"port", s.cfg.Server.Port,
			"endpoint", s.cfg.Server.MCPEndpoint,
		)
	default:
		s.mcp = mcp_golang.NewServer(stdio.NewStdioServerTransport())
		s.logger.Info("MCP stdio transport configured")
	}
	return nil
}

func (s *Server) Serve(ctx context.Context) error {
	s.logger.Info("starting MCP server")

	if s.cfg.Server.Transport == "http" {
		s.logger.Info("listening for HTTP MCP connections",
			"addr", fmt.Sprintf(":%d", s.cfg.Server.Port),
		)
	} else {
		s.logger.Info("listening for stdio MCP connections")
	}

	err := s.mcp.Serve()
	if err != nil {
		return fmt.Errorf("mcp serve: %w", err)
	}

	return nil
}
