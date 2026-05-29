package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	mcp_golang "github.com/metoro-io/mcp-golang"

	"mcp-susemanager/internal/api"
)

type MCPServices struct {
	System           *api.SystemService
	Audit            *api.AuditService
	Errata           *api.ErrataService
	Channel          *api.ChannelService
	ContentLifecycle *api.ContentLifecycleService
	Client           *api.Client
}

func NewServices(apiClient *api.Client) *MCPServices {
	return &MCPServices{
		System:           api.NewSystemService(apiClient),
		Audit:            api.NewAuditService(apiClient),
		Errata:           api.NewErrataService(apiClient),
		Channel:          api.NewChannelService(apiClient),
		ContentLifecycle: api.NewContentLifecycleService(apiClient),
		Client:           apiClient,
	}
}

func (s *Server) registerTools() error {
	svc := NewServices(s.api)

	tools := []struct {
		name        string
		description string
		handler     interface{}
	}{
		// ===== AUTH =====
		{
			name:        "suse_connect",
			description: "Test connection to SUSE Manager server and return API version",
			handler:     s.handleConnect(svc),
		},
		{
			name:        "suse_get_version",
			description: "Get the SUSE Manager API version",
			handler:     s.handleGetVersion(svc),
		},
		{
			name:        "suse_check_permissions",
			description: "Check current user permissions by listing accessible systems",
			handler:     s.handleCheckPermissions(svc),
		},
		{
			name:        "suse_disconnect",
			description: "Explicitly logout from SUSE Manager session",
			handler:     s.handleDisconnect(svc),
		},

		// ===== SYSTEMS =====
		{
			name:        "suse_list_systems",
			description: "List all registered systems in SUSE Manager",
			handler:     s.handleListSystems(svc),
		},
		{
			name:        "suse_get_system_details",
			description: "Get detailed information about a specific system by its ID",
			handler:     s.handleGetSystemDetails(svc),
		},
		{
			name:        "suse_search_system_by_hostname",
			description: "Search for systems by hostname using regex pattern",
			handler:     s.handleSearchSystemByHostname(svc),
		},
		{
			name:        "suse_get_system_hardware",
			description: "Get hardware inventory for a system (CPU, memory, DMI, devices)",
			handler:     s.handleGetSystemHardware(svc),
		},
		{
			name:        "suse_get_system_software",
			description: "Get list of installed packages on a system",
			handler:     s.handleGetSystemSoftware(svc),
		},
		{
			name:        "suse_get_system_events",
			description: "Get event history for a system with optional pagination",
			handler:     s.handleGetSystemEvents(svc),
		},
		{
			name:        "suse_schedule_reboot",
			description: "Schedule a reboot for a system",
			handler:     s.handleScheduleReboot(svc),
		},
		{
			name:        "suse_schedule_package_update",
			description: "Schedule a package update for a system",
			handler:     s.handleSchedulePackageUpdate(svc),
		},
		{
			name:        "suse_schedule_highstate",
			description: "Schedule a Salt highstate execution on a system",
			handler:     s.handleScheduleHighstate(svc),
		},
		{
			name:        "suse_apply_states",
			description: "Apply specific Salt states to a system",
			handler:     s.handleApplyStates(svc),
		},

		// ===== CVE / SECURITY =====
		{
			name:        "suse_list_cve_systems",
			description: "List systems by their patch status for a specific CVE identifier",
			handler:     s.handleListCveSystems(svc),
		},
		{
			name:        "suse_get_system_errata",
			description: "Get all relevant errata for a system",
			handler:     s.handleGetSystemErrata(svc),
		},
		{
			name:        "suse_schedule_errata",
			description: "Schedule errata (patches) to be applied to one or more systems",
			handler:     s.handleScheduleErrata(svc),
		},
		{
			name:        "suse_get_unscheduled_errata",
			description: "Get unscheduled errata for a system",
			handler:     s.handleGetUnscheduledErrata(svc),
		},

		// ===== CHANNELS =====
		{
			name:        "suse_list_channels",
			description: "List all software channels in SUSE Manager",
			handler:     s.handleListChannels(svc),
		},
		{
			name:        "suse_get_channel_details",
			description: "Get details about a specific software channel by its label",
			handler:     s.handleGetChannelDetails(svc),
		},
		{
			name:        "suse_list_channel_packages",
			description: "List packages in a software channel",
			handler:     s.handleListChannelPackages(svc),
		},
		{
			name:        "suse_list_system_channels",
			description: "List subscribed child channels for a system",
			handler:     s.handleListSystemChannels(svc),
		},
		{
			name:        "suse_change_system_channels",
			description: "Change the subscribed channels for a system (base + child channels)",
			handler:     s.handleChangeSystemChannels(svc),
		},
		{
			name:        "suse_create_repo",
			description: "Create a new repository in SUSE Manager",
			handler:     s.handleCreateRepo(svc),
		},
		{
			name:        "suse_create_channel",
			description: "Create a new software channel in SUSE Manager",
			handler:     s.handleCreateChannel(svc),
		},
		{
			name:        "suse_associate_repo_to_channel",
			description: "Associate a repository with a software channel",
			handler:     s.handleAssociateRepoToChannel(svc),
		},
		{
			name:        "suse_list_arches",
			description: "List available software channel architectures",
			handler:     s.handleListArches(svc),
		},
		{
			name:        "suse_list_repos",
			description: "List all user repositories in SUSE Manager",
			handler:     s.handleListRepos(svc),
		},
		{
			name:        "suse_list_channel_repos",
			description: "List repositories associated with a specific software channel",
			handler:     s.handleListChannelRepos(svc),
		},
		{
			name:        "suse_delete_channel",
			description: "Delete a custom software channel by its label",
			handler:     s.handleDeleteChannel(svc),
		},
		{
			name:        "suse_remove_repo",
			description: "Remove a repository by its label",
			handler:     s.handleRemoveRepo(svc),
		},

		// ===== CONTENT LIFECYCLE =====
		{
			name:        "suse_list_projects",
			description: "List all Content Lifecycle projects",
			handler:     s.handleListProjects(svc),
		},
		{
			name:        "suse_lookup_project",
			description: "Get details of a Content Lifecycle project by its label",
			handler:     s.handleLookupProject(svc),
		},
		{
			name:        "suse_create_project",
			description: "Create a new Content Lifecycle project",
			handler:     s.handleCreateProject(svc),
		},
		{
			name:        "suse_remove_project",
			description: "Remove a Content Lifecycle project by its label",
			handler:     s.handleRemoveProject(svc),
		},
		{
			name:        "suse_build_project",
			description: "Build a Content Lifecycle project",
			handler:     s.handleBuildProject(svc),
		},
		{
			name:        "suse_promote_project",
			description: "Promote a Content Lifecycle project environment",
			handler:     s.handlePromoteProject(svc),
		},
		{
			name:        "suse_list_project_environments",
			description: "List environments in a Content Lifecycle project",
			handler:     s.handleListProjectEnvironments(svc),
		},
		{
			name:        "suse_create_environment",
			description: "Create a new environment in a Content Lifecycle project",
			handler:     s.handleCreateEnvironment(svc),
		},
		{
			name:        "suse_attach_source",
			description: "Attach a source (channel) to a Content Lifecycle project",
			handler:     s.handleAttachSource(svc),
		},
		{
			name:        "suse_detach_source",
			description: "Detach a source (channel) from a Content Lifecycle project",
			handler:     s.handleDetachSource(svc),
		},
		{
			name:        "suse_list_project_sources",
			description: "List sources attached to a Content Lifecycle project (optional: filter by source type)",
			handler:     s.handleListProjectSources(svc),
		},
		{
			name:        "suse_set_sources",
			description: "Replace all sources of a given type in a Content Lifecycle project",
			handler:     s.handleSetSources(svc),
		},
	}

	for _, tool := range tools {
		if err := s.mcp.RegisterTool(tool.name, tool.description, tool.handler); err != nil {
			return fmt.Errorf("register tool %s: %w", tool.name, err)
		}
		s.logger.Debug("tool registered", "name", tool.name)
	}

	s.logger.Info("tools registered", "count", len(tools))
	return nil
}

// ----- Tool argument types -----

type EmptyArgs struct{}

type SystemIDArgs struct {
	SystemID int `json:"system_id" jsonschema:"required,description=System ID in SUSE Manager"`
}

type SystemIDsArgs struct {
	SystemIDs []int  `json:"system_ids" jsonschema:"required,description=List of system IDs"`
	Earliest  string `json:"earliest,omitempty" jsonschema:"description=Earliest occurrence time (ISO 8601)"`
}

type SearchNameArgs struct {
	Pattern string `json:"pattern" jsonschema:"required,description=Regex pattern to search hostnames"`
}

type PaginationArgs struct {
	SystemID int `json:"system_id" jsonschema:"required,description=System ID"`
	Offset   int `json:"offset,omitempty" jsonschema:"description=Pagination offset"`
	Limit    int `json:"limit,omitempty" jsonschema:"description=Maximum number of results"`
}

type ScheduleActionArgs struct {
	SystemID int    `json:"system_id" jsonschema:"required,description=System ID"`
	Earliest string `json:"earliest,omitempty" jsonschema:"description=Earliest occurrence time (ISO 8601, e.g. 2025-01-01T00:00:00)"`
}

type ScheduleHighstateArgs struct {
	SystemID int    `json:"system_id" jsonschema:"required,description=System ID"`
	Test     bool   `json:"test,omitempty" jsonschema:"description=Test mode (dry run)"`
	Earliest string `json:"earliest,omitempty" jsonschema:"description=Earliest occurrence time (ISO 8601)"`
}

type ApplyStatesArgs struct {
	SystemID   int      `json:"system_id" jsonschema:"required,description=System ID"`
	StateNames []string `json:"state_names" jsonschema:"required,description=List of Salt state names to apply"`
	Test       bool     `json:"test,omitempty" jsonschema:"description=Test mode (dry run)"`
	Earliest   string   `json:"earliest,omitempty" jsonschema:"description=Earliest occurrence time (ISO 8601)"`
}

type CveArgs struct {
	CveIdentifier string   `json:"cve_identifier" jsonschema:"required,description=CVE identifier (e.g. CVE-2025-1234)"`
	PatchStatuses []string `json:"patch_statuses,omitempty" jsonschema:"description=Filter by patch statuses: AFFECTED_PATCH_INAPPLICABLE, AFFECTED_PATCH_APPLICABLE, NOT_AFFECTED, PATCHED"`
}

type ErrataScheduleArgs struct {
	SystemIDs []int  `json:"system_ids" jsonschema:"required,description=List of system IDs"`
	ErrataIDs []int  `json:"errata_ids" jsonschema:"required,description=List of errata IDs to apply"`
	Earliest  string `json:"earliest,omitempty" jsonschema:"description=Earliest occurrence time (ISO 8601)"`
}

type ChannelLabelArgs struct {
	Label string `json:"label" jsonschema:"required,description=Software channel label"`
}

type SystemChannelChangeArgs struct {
	SystemID        int      `json:"system_id" jsonschema:"required,description=System ID"`
	BaseChannel     string   `json:"base_channel,omitempty" jsonschema:"description=Base channel label"`
	ChildChannelIDs []int    `json:"child_channel_ids,omitempty" jsonschema:"description=Child channel IDs"`
	Earliest        string   `json:"earliest,omitempty" jsonschema:"description=Earliest occurrence time (ISO 8601)"`
}

// ----- Tool Handler Implementations -----

func (s *Server) handleConnect(svc *MCPServices) interface{} {
	return func(args EmptyArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		_, version, err := svc.Client.TestConnection(ctx)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Connection failed: %v", err)),
			), nil
		}
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(fmt.Sprintf("Connected to SUSE Manager (API version: %s)", version)),
		), nil
	}
}

func (s *Server) handleGetVersion(svc *MCPServices) interface{} {
	return func(args EmptyArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		version, err := svc.Client.GetVersion(ctx)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error getting version: %v", err)),
			), nil
		}
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(fmt.Sprintf("SUSE Manager API version: %s", version.Version)),
		), nil
	}
}

func (s *Server) handleCheckPermissions(svc *MCPServices) interface{} {
	return func(args EmptyArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		systems, err := svc.System.ListSystems(ctx)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error checking permissions: %v", err)),
			), nil
		}
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(fmt.Sprintf("User has access to %d systems. Session key: %s", len(systems), maskKey(svc.Client.SessionKey()))),
		), nil
	}
}

func (s *Server) handleDisconnect(svc *MCPServices) interface{} {
	return func(args EmptyArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := svc.Client.Logout(ctx); err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Logout warning: %v", err)),
			), nil
		}
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent("Logged out successfully"),
		), nil
	}
}

func (s *Server) handleListSystems(svc *MCPServices) interface{} {
	return func(args EmptyArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		systems, err := svc.System.ListSystems(ctx)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error listing systems: %v", err)),
			), nil
		}

		data, _ := json.MarshalIndent(systems, "", "  ")
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(string(data)),
		), nil
	}
}

func (s *Server) handleGetSystemDetails(svc *MCPServices) interface{} {
	return func(args SystemIDArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		details, err := svc.System.GetDetails(ctx, args.SystemID)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error getting system details: %v", err)),
			), nil
		}

		data, _ := json.MarshalIndent(details, "", "  ")
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(string(data)),
		), nil
	}
}

func (s *Server) handleSearchSystemByHostname(svc *MCPServices) interface{} {
	return func(args SearchNameArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		systems, err := svc.System.SearchByName(ctx, args.Pattern)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error searching systems: %v", err)),
			), nil
		}

		data, _ := json.MarshalIndent(systems, "", "  ")
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(string(data)),
		), nil
	}
}

func (s *Server) handleGetSystemHardware(svc *MCPServices) interface{} {
	return func(args SystemIDArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		cpu, cpuErr := svc.System.GetCPURaw(ctx, args.SystemID)
		mem, memErr := svc.System.GetMemoryRaw(ctx, args.SystemID)
		dmi, dmiErr := svc.System.GetDmiRaw(ctx, args.SystemID)
		devices, devErr := svc.System.GetDevicesRaw(ctx, args.SystemID)

		result := map[string]interface{}{
			"system_id": args.SystemID,
		}

		if cpuErr != nil {
			result["cpu_error"] = cpuErr.Error()
		} else {
			result["cpu"] = cpu
		}
		if memErr != nil {
			result["memory_error"] = memErr.Error()
		} else {
			result["memory"] = mem
		}
		if dmiErr != nil {
			result["dmi_error"] = dmiErr.Error()
		} else {
			result["dmi"] = dmi
		}
		if devErr != nil {
			result["devices_error"] = devErr.Error()
		} else {
			result["devices"] = devices
		}

		data, _ := json.MarshalIndent(result, "", "  ")
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(string(data)),
		), nil
	}
}

func (s *Server) handleGetSystemSoftware(svc *MCPServices) interface{} {
	return func(args SystemIDArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		packages, err := svc.System.ListInstalledPackages(ctx, args.SystemID)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error listing installed packages: %v", err)),
			), nil
		}

		data, _ := json.MarshalIndent(packages, "", "  ")
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(string(data)),
		), nil
	}
}

func (s *Server) handleGetSystemEvents(svc *MCPServices) interface{} {
	return func(args PaginationArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		events, err := svc.System.GetEventHistory(ctx, args.SystemID, args.Offset, args.Limit)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error getting event history: %v", err)),
			), nil
		}

		data, _ := json.MarshalIndent(events, "", "  ")
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(string(data)),
		), nil
	}
}

func (s *Server) handleScheduleReboot(svc *MCPServices) interface{} {
	return func(args ScheduleActionArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		actionID, err := svc.System.ScheduleReboot(ctx, args.SystemID, args.Earliest)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error scheduling reboot: %v", err)),
			), nil
		}
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(fmt.Sprintf("Reboot scheduled successfully. Action ID: %d", actionID)),
		), nil
	}
}

func (s *Server) handleSchedulePackageUpdate(svc *MCPServices) interface{} {
	return func(args ScheduleActionArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		actionID, err := svc.System.SchedulePackageUpdate(ctx, args.SystemID, args.Earliest)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error scheduling package update: %v", err)),
			), nil
		}
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(fmt.Sprintf("Package update scheduled successfully. Action ID: %d", actionID)),
		), nil
	}
}

func (s *Server) handleScheduleHighstate(svc *MCPServices) interface{} {
	return func(args ScheduleHighstateArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		actionID, err := svc.System.ScheduleApplyHighstate(ctx, args.SystemID, args.Earliest, args.Test)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error scheduling highstate: %v", err)),
			), nil
		}
		msg := fmt.Sprintf("Highstate scheduled successfully. Action ID: %d", actionID)
		if args.Test {
			msg += " (test mode)"
		}
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(msg),
		), nil
	}
}

func (s *Server) handleApplyStates(svc *MCPServices) interface{} {
	return func(args ApplyStatesArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		actionID, err := svc.System.ScheduleApplyStates(ctx, args.SystemID, args.StateNames, args.Earliest, args.Test)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error applying states: %v", err)),
			), nil
		}
		msg := fmt.Sprintf("States applied successfully. Action ID: %d", actionID)
		if args.Test {
			msg += " (test mode)"
		}
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(msg),
		), nil
	}
}

func (s *Server) handleListCveSystems(svc *MCPServices) interface{} {
	return func(args CveArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		systems, err := svc.Audit.ListSystemsByPatchStatus(ctx, args.CveIdentifier, args.PatchStatuses)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error listing CVE systems: %v", err)),
			), nil
		}

		data, _ := json.MarshalIndent(systems, "", "  ")
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(string(data)),
		), nil
	}
}

func (s *Server) handleGetSystemErrata(svc *MCPServices) interface{} {
	return func(args SystemIDArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		errata, err := svc.Errata.GetRelevantErrata(ctx, args.SystemID)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error getting errata: %v", err)),
			), nil
		}

		data, _ := json.MarshalIndent(errata, "", "  ")
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(string(data)),
		), nil
	}
}

func (s *Server) handleScheduleErrata(svc *MCPServices) interface{} {
	return func(args ErrataScheduleArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		actionID, err := svc.Errata.ScheduleApplyErrata(ctx, args.SystemIDs, args.ErrataIDs, args.Earliest)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error scheduling errata: %v", err)),
			), nil
		}
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(fmt.Sprintf("Errata scheduled successfully. Action ID: %d", actionID)),
		), nil
	}
}

func (s *Server) handleGetUnscheduledErrata(svc *MCPServices) interface{} {
	return func(args SystemIDArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		errata, err := svc.Errata.GetUnscheduledErrata(ctx, args.SystemID)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error getting unscheduled errata: %v", err)),
			), nil
		}

		data, _ := json.MarshalIndent(errata, "", "  ")
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(string(data)),
		), nil
	}
}

func (s *Server) handleListChannels(svc *MCPServices) interface{} {
	return func(args EmptyArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		channels, err := svc.Channel.ListAllChannels(ctx)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error listing channels: %v", err)),
			), nil
		}

		data, _ := json.MarshalIndent(channels, "", "  ")
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(string(data)),
		), nil
	}
}

func (s *Server) handleGetChannelDetails(svc *MCPServices) interface{} {
	return func(args ChannelLabelArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		details, err := svc.Channel.GetDetails(ctx, args.Label)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error getting channel details: %v", err)),
			), nil
		}

		data, _ := json.MarshalIndent(details, "", "  ")
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(string(data)),
		), nil
	}
}

func (s *Server) handleListChannelPackages(svc *MCPServices) interface{} {
	return func(args ChannelLabelArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		packages, err := svc.Channel.ListPackages(ctx, args.Label)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error listing channel packages: %v", err)),
			), nil
		}

		data, _ := json.MarshalIndent(packages, "", "  ")
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(string(data)),
		), nil
	}
}

func (s *Server) handleListSystemChannels(svc *MCPServices) interface{} {
	return func(args SystemIDArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		channels, err := svc.Channel.ListSubscribedChildChannels(ctx, args.SystemID)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error listing system channels: %v", err)),
			), nil
		}

		data, _ := json.MarshalIndent(channels, "", "  ")
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(string(data)),
		), nil
	}
}

func (s *Server) handleChangeSystemChannels(svc *MCPServices) interface{} {
	return func(args SystemChannelChangeArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if args.BaseChannel != "" {
			actionID, err := svc.Channel.ScheduleChangeChannels(ctx, []int{args.SystemID}, args.BaseChannel, nil, args.Earliest)
			if err != nil {
				return mcp_golang.NewToolResponse(
					mcp_golang.NewTextContent(fmt.Sprintf("Error changing channels: %v", err)),
				), nil
			}
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Channel change scheduled. Action ID: %d", actionID)),
			), nil
		}
		if args.ChildChannelIDs != nil {
			result, err := svc.Channel.SetChildChannels(ctx, args.SystemID, args.ChildChannelIDs)
			if err != nil {
				return mcp_golang.NewToolResponse(
					mcp_golang.NewTextContent(fmt.Sprintf("Error setting child channels: %v", err)),
				), nil
			}
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Child channels updated. Result: %d", result)),
			), nil
		}
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent("No channel changes specified. Provide base_channel or child_channel_ids."),
		), nil
	}
}

type CreateRepoArgs struct {
	Label    string `json:"label" jsonschema:"required,description=Repository label"`
	RepoType string `json:"repo_type" jsonschema:"required,description=Repository type (e.g. yum, uln)"`
	URL      string `json:"url" jsonschema:"required,description=Repository URL"`
}

func (s *Server) handleCreateRepo(svc *MCPServices) interface{} {
	return func(args CreateRepoArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		repo, err := svc.Channel.CreateRepo(ctx, args.Label, args.RepoType, args.URL)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error creating repository: %v", err)),
			), nil
		}

		data, _ := json.MarshalIndent(repo, "", "  ")
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(string(data)),
		), nil
	}
}

type CreateChannelArgs struct {
	Label       string `json:"label" jsonschema:"required,description=Channel label (unique identifier)"`
	Name        string `json:"name" jsonschema:"required,description=Channel display name"`
	Summary     string `json:"summary" jsonschema:"required,description=Short channel summary"`
	ArchLabel   string `json:"arch_label" jsonschema:"required,description=Architecture label (e.g. x86_64, i386). Use suse_list_arches to see available options"`
	ParentLabel string `json:"parent_label,omitempty" jsonschema:"description=Parent channel label (empty string for base channel)"`
}

type AssociateRepoArgs struct {
	ChannelLabel string `json:"channel_label" jsonschema:"required,description=Software channel label"`
	RepoLabel    string `json:"repo_label" jsonschema:"required,description=Repository label to associate"`
}

func (s *Server) handleCreateChannel(svc *MCPServices) interface{} {
	return func(args CreateChannelArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		result, err := svc.Channel.Create(ctx, args.Label, args.Name, args.Summary, args.ArchLabel, args.ParentLabel)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error creating channel: %v", err)),
			), nil
		}
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(fmt.Sprintf("Channel created successfully. Result: %d", result)),
		), nil
	}
}

func (s *Server) handleAssociateRepoToChannel(svc *MCPServices) interface{} {
	return func(args AssociateRepoArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		channel, err := svc.Channel.AssociateRepo(ctx, args.ChannelLabel, args.RepoLabel)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error associating repo to channel: %v", err)),
			), nil
		}

		data, _ := json.MarshalIndent(channel, "", "  ")
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(string(data)),
		), nil
	}
}

func (s *Server) handleListArches(svc *MCPServices) interface{} {
	return func(args EmptyArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		arches, err := svc.Channel.ListArches(ctx)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error listing arches: %v", err)),
			), nil
		}

		data, _ := json.MarshalIndent(arches, "", "  ")
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(string(data)),
		), nil
	}
}

type ChannelLabelOnlyArgs struct {
	ChannelLabel string `json:"channel_label" jsonschema:"required,description=Software channel label"`
}

type RepoLabelArgs struct {
	RepoLabel string `json:"repo_label" jsonschema:"required,description=Repository label"`
}

type ChannelReposArgs struct {
	ChannelLabel string `json:"channel_label" jsonschema:"required,description=Software channel label"`
}

func (s *Server) handleListRepos(svc *MCPServices) interface{} {
	return func(args EmptyArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		repos, err := svc.Channel.ListUserRepos(ctx)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error listing repos: %v", err)),
			), nil
		}

		data, _ := json.MarshalIndent(repos, "", "  ")
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(string(data)),
		), nil
	}
}

func (s *Server) handleListChannelRepos(svc *MCPServices) interface{} {
	return func(args ChannelReposArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		repos, err := svc.Channel.ListChannelRepos(ctx, args.ChannelLabel)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error listing channel repos: %v", err)),
			), nil
		}

		data, _ := json.MarshalIndent(repos, "", "  ")
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(string(data)),
		), nil
	}
}

func (s *Server) handleDeleteChannel(svc *MCPServices) interface{} {
	return func(args ChannelLabelOnlyArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		result, err := svc.Channel.DeleteChannel(ctx, args.ChannelLabel)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error deleting channel: %v", err)),
			), nil
		}
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(fmt.Sprintf("Channel deleted successfully. Result: %d", result)),
		), nil
	}
}

func (s *Server) handleRemoveRepo(svc *MCPServices) interface{} {
	return func(args RepoLabelArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		result, err := svc.Channel.RemoveRepo(ctx, args.RepoLabel)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error removing repo: %v", err)),
			), nil
		}
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(fmt.Sprintf("Repository removed successfully. Result: %d", result)),
		), nil
	}
}

// ===== CONTENT LIFECYCLE HANDLERS =====

type CreateProjectArgs struct {
	Label       string `json:"label" jsonschema:"required,description=Project label (unique identifier)"`
	Name        string `json:"name" jsonschema:"required,description=Project display name"`
	Description string `json:"description" jsonschema:"required,description=Project description"`
}

type ProjectLabelArgs struct {
	ProjectLabel string `json:"project_label" jsonschema:"required,description=Content Lifecycle project label"`
}

type BuildProjectArgs struct {
	ProjectLabel string `json:"project_label" jsonschema:"required,description=Content Lifecycle project label"`
	Message      string `json:"message,omitempty" jsonschema:"description=Optional log message for the build"`
}

type PromoteProjectArgs struct {
	ProjectLabel   string `json:"project_label" jsonschema:"required,description=Content Lifecycle project label"`
	EnvironmentLabel string `json:"environment_label" jsonschema:"required,description=Environment label to promote"`
}

func (s *Server) handleListProjects(svc *MCPServices) interface{} {
	return func(args EmptyArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		projects, err := svc.ContentLifecycle.ListProjects(ctx)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error listing projects: %v", err)),
			), nil
		}

		data, _ := json.MarshalIndent(projects, "", "  ")
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(string(data)),
		), nil
	}
}

func (s *Server) handleLookupProject(svc *MCPServices) interface{} {
	return func(args ProjectLabelArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		project, err := svc.ContentLifecycle.LookupProject(ctx, args.ProjectLabel)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error looking up project: %v", err)),
			), nil
		}

		data, _ := json.MarshalIndent(project, "", "  ")
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(string(data)),
		), nil
	}
}

func (s *Server) handleCreateProject(svc *MCPServices) interface{} {
	return func(args CreateProjectArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		project, err := svc.ContentLifecycle.CreateProject(ctx, args.Label, args.Name, args.Description)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error creating project: %v", err)),
			), nil
		}

		data, _ := json.MarshalIndent(project, "", "  ")
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(string(data)),
		), nil
	}
}

func (s *Server) handleRemoveProject(svc *MCPServices) interface{} {
	return func(args ProjectLabelArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		result, err := svc.ContentLifecycle.RemoveProject(ctx, args.ProjectLabel)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error removing project: %v", err)),
			), nil
		}
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(fmt.Sprintf("Project removed successfully. Result: %d", result)),
		), nil
	}
}

func (s *Server) handleBuildProject(svc *MCPServices) interface{} {
	return func(args BuildProjectArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		result, err := svc.ContentLifecycle.BuildProject(ctx, args.ProjectLabel, args.Message)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error building project: %v", err)),
			), nil
		}
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(fmt.Sprintf("Project build started. Result: %d", result)),
		), nil
	}
}

func (s *Server) handlePromoteProject(svc *MCPServices) interface{} {
	return func(args PromoteProjectArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		result, err := svc.ContentLifecycle.PromoteProject(ctx, args.ProjectLabel, args.EnvironmentLabel)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error promoting project: %v", err)),
			), nil
		}
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(fmt.Sprintf("Project promoted successfully. Result: %d", result)),
		), nil
	}
}

func (s *Server) handleListProjectEnvironments(svc *MCPServices) interface{} {
	return func(args ProjectLabelArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		envs, err := svc.ContentLifecycle.ListProjectEnvironments(ctx, args.ProjectLabel)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error listing environments: %v", err)),
			), nil
		}

		data, _ := json.MarshalIndent(envs, "", "  ")
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(string(data)),
		), nil
	}
}

type CreateEnvironmentArgs struct {
	ProjectLabel     string `json:"project_label" jsonschema:"required,description=Content Lifecycle project label"`
	PredecessorLabel string `json:"predecessor_label" jsonschema:"required,description=Predecessor environment label (use empty string for first environment)"`
	EnvLabel         string `json:"env_label" jsonschema:"required,description=New environment label"`
	Name             string `json:"name" jsonschema:"required,description=Environment display name"`
	Description      string `json:"description" jsonschema:"required,description=Environment description"`
}

func (s *Server) handleCreateEnvironment(svc *MCPServices) interface{} {
	return func(args CreateEnvironmentArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		env, err := svc.ContentLifecycle.CreateEnvironment(ctx, args.ProjectLabel, args.PredecessorLabel, args.EnvLabel, args.Name, args.Description)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error creating environment: %v", err)),
			), nil
		}

		data, _ := json.MarshalIndent(env, "", "  ")
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(string(data)),
		), nil
	}
}

type AttachSourceArgs struct {
	ProjectLabel string `json:"project_label" jsonschema:"required,description=Content Lifecycle project label"`
	SourceType   string `json:"source_type" jsonschema:"required,description=Source type: software, state, or image"`
	SourceLabel  string `json:"source_label" jsonschema:"required,description=Source (channel) label to attach"`
	Position     int    `json:"position,omitempty" jsonschema:"description=Optional position in the source list"`
}

type DetachSourceArgs struct {
	ProjectLabel string `json:"project_label" jsonschema:"required,description=Content Lifecycle project label"`
	SourceType   string `json:"source_type" jsonschema:"required,description=Source type: software, state, or image"`
	SourceLabel  string `json:"source_label" jsonschema:"required,description=Source (channel) label to detach"`
}

type ListProjectSourcesArgs struct {
	ProjectLabel string `json:"project_label" jsonschema:"required,description=Content Lifecycle project label"`
	SourceType   string `json:"source_type,omitempty" jsonschema:"description=Optional source type filter: software, state, or image"`
}

type SetSourcesArgs struct {
	ProjectLabel string   `json:"project_label" jsonschema:"required,description=Content Lifecycle project label"`
	SourceType   string   `json:"source_type" jsonschema:"required,description=Source type: software, state, or image"`
	SourceLabels []string `json:"source_labels" jsonschema:"required,description=List of source (channel) labels to set"`
}

func (s *Server) handleAttachSource(svc *MCPServices) interface{} {
	return func(args AttachSourceArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		result, err := svc.ContentLifecycle.AttachSource(ctx, args.ProjectLabel, args.SourceType, args.SourceLabel, args.Position)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error attaching source: %v", err)),
			), nil
		}

		data, _ := json.MarshalIndent(result, "", "  ")
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(string(data)),
		), nil
	}
}

func (s *Server) handleDetachSource(svc *MCPServices) interface{} {
	return func(args DetachSourceArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		result, err := svc.ContentLifecycle.DetachSource(ctx, args.ProjectLabel, args.SourceType, args.SourceLabel)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error detaching source: %v", err)),
			), nil
		}

		data, _ := json.MarshalIndent(result, "", "  ")
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(string(data)),
		), nil
	}
}

func (s *Server) handleListProjectSources(svc *MCPServices) interface{} {
	return func(args ListProjectSourcesArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		sources, err := svc.ContentLifecycle.ListProjectSources(ctx, args.ProjectLabel, args.SourceType)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error listing sources: %v", err)),
			), nil
		}

		data, _ := json.MarshalIndent(sources, "", "  ")
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(string(data)),
		), nil
	}
}

func (s *Server) handleSetSources(svc *MCPServices) interface{} {
	return func(args SetSourcesArgs) (*mcp_golang.ToolResponse, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		result, err := svc.ContentLifecycle.SetSources(ctx, args.ProjectLabel, args.SourceType, args.SourceLabels)
		if err != nil {
			return mcp_golang.NewToolResponse(
				mcp_golang.NewTextContent(fmt.Sprintf("Error setting sources: %v", err)),
			), nil
		}

		data, _ := json.MarshalIndent(result, "", "  ")
		return mcp_golang.NewToolResponse(
			mcp_golang.NewTextContent(string(data)),
		), nil
	}
}

func maskKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "..." + key[len(key)-4:]
}
