package mcp

import (
	"fmt"

	mcp_golang "github.com/metoro-io/mcp-golang"
)

type SecurityAuditArgs struct {
	Target string `json:"target,omitempty" jsonschema:"description=Optional: focus on a specific CVE or system"`
}

type PatchStatusArgs struct {
	Group string `json:"group,omitempty" jsonschema:"description=Optional: filter by system group name"`
}

type SystemOverviewArgs struct {
	Hostname string `json:"hostname,omitempty" jsonschema:"description=Optional: focus on a specific hostname pattern"`
}

func (s *Server) registerPrompts() error {
	prompts := []struct {
		name        string
		description string
		handler     interface{}
	}{
		{
			name:        "security-audit",
			description: "Audit security: list unfixed CVEs and affected systems. Provide a structured security assessment.",
			handler: s.securityAuditPrompt(),
		},
		{
			name:        "patch-status",
			description: "Get a summary of patch status across all systems or a specific group.",
			handler: s.patchStatusPrompt(),
		},
		{
			name:        "system-overview",
			description: "Get an overview of all systems, their hardware, and subscribed channels.",
			handler: s.systemOverviewPrompt(),
		},
	}

	for _, p := range prompts {
		if err := s.mcp.RegisterPrompt(p.name, p.description, p.handler); err != nil {
			return fmt.Errorf("register prompt %s: %w", p.name, err)
		}
		s.logger.Debug("prompt registered", "name", p.name)
	}

	s.logger.Info("prompts registered", "count", len(prompts))
	return nil
}

func (s *Server) securityAuditPrompt() interface{} {
	return func(args SecurityAuditArgs) (*mcp_golang.PromptResponse, error) {
		msg := "You are a security auditor for SUSE Manager. "
		msg += "Analyze the following data and provide a structured security assessment:\n\n"
		msg += "1. List all systems and their CVE patch status\n"
		if args.Target != "" {
			msg += fmt.Sprintf("2. Focus specifically on: %s\n", args.Target)
		}
		msg += "3. Identify systems that are AFFECTED_PATCH_INAPPLICABLE (affected but patch in unassigned channel)\n"
		msg += "4. Identify systems that are AFFECTED_PATCH_APPLICABLE (affected but patch not yet applied)\n"
		msg += "5. Provide prioritized remediation recommendations\n"
		msg += "\nUse the suse_list_cve_systems and suse_get_system_errata tools to gather the necessary data."

		return mcp_golang.NewPromptResponse(
			"security-audit",
			mcp_golang.NewPromptMessage(
				mcp_golang.NewTextContent(msg),
				mcp_golang.RoleUser,
			),
		), nil
	}
}

func (s *Server) patchStatusPrompt() interface{} {
	return func(args PatchStatusArgs) (*mcp_golang.PromptResponse, error) {
		msg := "You are a patch management assistant for SUSE Manager. "
		msg += "Analyze the patch status and provide a summary:\n\n"
		msg += "1. List all systems and their patch status\n"
		msg += "2. Identify systems that are out of date\n"
		msg += "3. Identify systems that need reboots\n"
		msg += "4. Check for recent event history\n"
		if args.Group != "" {
			msg += fmt.Sprintf("5. Filter by system group: %s\n", args.Group)
		}
		msg += "6. Provide a prioritized patching plan\n"
		msg += "\nUse suse_list_systems, suse_get_system_errata, and suse_get_system_events."

		return mcp_golang.NewPromptResponse(
			"patch-status",
			mcp_golang.NewPromptMessage(
				mcp_golang.NewTextContent(msg),
				mcp_golang.RoleUser,
			),
		), nil
	}
}

func (s *Server) systemOverviewPrompt() interface{} {
	return func(args SystemOverviewArgs) (*mcp_golang.PromptResponse, error) {
		msg := "You are a system administrator for SUSE Manager. "
		msg += "Provide a comprehensive overview:\n\n"
		msg += "1. List all registered systems\n"
		msg += "2. For each system, show:\n"
		msg += "   - System ID, name, hostname\n"
		msg += "   - Hardware info (CPU, memory, DMI)\n"
		msg += "   - Installed packages count\n"
		msg += "   - Subscribed channels\n"
		msg += "   - Recent events\n"
		if args.Hostname != "" {
			msg += fmt.Sprintf("3. Focus on: %s\n", args.Hostname)
		}
		msg += "\nUse suse_list_systems, suse_get_system_details, and suse_get_system_hardware."

		return mcp_golang.NewPromptResponse(
			"system-overview",
			mcp_golang.NewPromptMessage(
				mcp_golang.NewTextContent(msg),
				mcp_golang.RoleUser,
			),
		), nil
	}
}
