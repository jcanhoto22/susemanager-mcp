package api

import (
	"context"
	"fmt"

	"mcp-susemanager/internal/models"
)

type SystemService struct {
	client func(ctx context.Context, method string, args []interface{}, reply interface{}) error
}

func NewSystemService(apiClient *Client) *SystemService {
	return &SystemService{
		client: apiClient.Call,
	}
}

func (s *SystemService) ListSystems(ctx context.Context) ([]models.System, error) {
	var result []models.System
	err := s.client(ctx, "system.listSystems", nil, &result)
	if err != nil {
		return nil, fmt.Errorf("list systems: %w", err)
	}
	return result, nil
}

func (s *SystemService) GetDetails(ctx context.Context, sid int) (*models.SystemDetails, error) {
	var result models.SystemDetails
	err := s.client(ctx, "system.getDetails", []interface{}{sid}, &result)
	if err != nil {
		return nil, fmt.Errorf("get system details: %w", err)
	}
	return &result, nil
}

func (s *SystemService) SearchByName(ctx context.Context, regex string) ([]models.System, error) {
	var result []models.System
	err := s.client(ctx, "system.searchByName", []interface{}{regex}, &result)
	if err != nil {
		return nil, fmt.Errorf("search system by name: %w", err)
	}
	return result, nil
}

func (s *SystemService) GetID(ctx context.Context, name string) (int, error) {
	var result int
	err := s.client(ctx, "system.getId", []interface{}{name}, &result)
	if err != nil {
		return 0, fmt.Errorf("get system id: %w", err)
	}
	return result, nil
}

func (s *SystemService) GetCPU(ctx context.Context, sid int) (*models.CPUInfo, error) {
	var result models.CPUInfo
	err := s.client(ctx, "system.getCpu", []interface{}{sid}, &result)
	if err != nil {
		return nil, fmt.Errorf("get cpu: %w", err)
	}
	return &result, nil
}

func (s *SystemService) GetMemory(ctx context.Context, sid int) (*models.MemoryInfo, error) {
	var result models.MemoryInfo
	err := s.client(ctx, "system.getMemory", []interface{}{sid}, &result)
	if err != nil {
		return nil, fmt.Errorf("get memory: %w", err)
	}
	return &result, nil
}

func (s *SystemService) GetDmi(ctx context.Context, sid int) (*models.DMIInfo, error) {
	var result models.DMIInfo
	err := s.client(ctx, "system.getDmi", []interface{}{sid}, &result)
	if err != nil {
		return nil, fmt.Errorf("get dmi: %w", err)
	}
	return &result, nil
}

func (s *SystemService) GetDevices(ctx context.Context, sid int) ([]map[string]string, error) {
	var result []map[string]string
	err := s.client(ctx, "system.getDevices", []interface{}{sid}, &result)
	if err != nil {
		return nil, fmt.Errorf("get devices: %w", err)
	}
	return result, nil
}

func (s *SystemService) ListInstalledPackages(ctx context.Context, sid int) ([]models.InstalledPackage, error) {
	var result []models.InstalledPackage
	err := s.client(ctx, "system.listInstalledPackages", []interface{}{sid}, &result)
	if err != nil {
		return nil, fmt.Errorf("list installed packages: %w", err)
	}
	return result, nil
}

func (s *SystemService) GetEventHistory(ctx context.Context, sid int, offset, limit int) ([]models.EventHistory, error) {
	var result []models.EventHistory
	if offset > 0 || limit > 0 {
		err := s.client(ctx, "system.getEventHistory", []interface{}{sid, offset, limit}, &result)
		if err != nil {
			return nil, fmt.Errorf("get event history: %w", err)
		}
	} else {
		err := s.client(ctx, "system.getEventHistory", []interface{}{sid}, &result)
		if err != nil {
			return nil, fmt.Errorf("get event history: %w", err)
		}
	}
	return result, nil
}

func (s *SystemService) ScheduleReboot(ctx context.Context, sid int, earliest string) (int, error) {
	var result int
	args := []interface{}{sid}
	if earliest != "" {
		args = append(args, earliest)
	}
	err := s.client(ctx, "system.scheduleReboot", args, &result)
	if err != nil {
		return 0, fmt.Errorf("schedule reboot: %w", err)
	}
	return result, nil
}

func (s *SystemService) SchedulePackageUpdate(ctx context.Context, sid int, earliest string) (int, error) {
	var result int
	args := []interface{}{sid}
	if earliest != "" {
		args = append(args, earliest)
	}
	err := s.client(ctx, "system.schedulePackageUpdate", args, &result)
	if err != nil {
		return 0, fmt.Errorf("schedule package update: %w", err)
	}
	return result, nil
}

func (s *SystemService) ScheduleApplyHighstate(ctx context.Context, sid int, earliest string, test bool) (int, error) {
	var result int
	args := []interface{}{sid}
	if earliest != "" {
		args = append(args, earliest, test)
	} else {
		args = append(args, test)
	}
	err := s.client(ctx, "system.scheduleApplyHighstate", []interface{}{sid, earliest, test}, &result)
	if err != nil {
		return 0, fmt.Errorf("schedule highstate: %w", err)
	}
	return result, nil
}

func (s *SystemService) ScheduleApplyStates(ctx context.Context, sid int, stateNames []string, earliest string, test bool) (int, error) {
	var result int
	args := []interface{}{sid, stateNames}
	if earliest != "" {
		args = append(args, earliest, test)
	} else {
		args = append(args, test)
	}
	err := s.client(ctx, "system.scheduleApplyStates", args, &result)
	if err != nil {
		return 0, fmt.Errorf("schedule apply states: %w", err)
	}
	return result, nil
}

func (s *SystemService) GetCPURaw(ctx context.Context, sid int) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := s.client(ctx, "system.getCpu", []interface{}{sid}, &result)
	if err != nil {
		return nil, fmt.Errorf("get cpu: %w", err)
	}
	return result, nil
}

func (s *SystemService) GetMemoryRaw(ctx context.Context, sid int) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := s.client(ctx, "system.getMemory", []interface{}{sid}, &result)
	if err != nil {
		return nil, fmt.Errorf("get memory: %w", err)
	}
	return result, nil
}

func (s *SystemService) GetDmiRaw(ctx context.Context, sid int) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := s.client(ctx, "system.getDmi", []interface{}{sid}, &result)
	if err != nil {
		return nil, fmt.Errorf("get dmi: %w", err)
	}
	return result, nil
}

func (s *SystemService) GetDevicesRaw(ctx context.Context, sid int) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	err := s.client(ctx, "system.getDevices", []interface{}{sid}, &result)
	if err != nil {
		return nil, fmt.Errorf("get devices: %w", err)
	}
	return result, nil
}

func (s *SystemService) GetPillar(ctx context.Context, minionID, category string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := s.client(ctx, "system.getPillar", []interface{}{minionID, category}, &result)
	if err != nil {
		return nil, fmt.Errorf("get pillar: %w", err)
	}
	return result, nil
}

func (s *SystemService) ListGroups(ctx context.Context, sid int) ([]models.SystemGroup, error) {
	var result []models.SystemGroup
	err := s.client(ctx, "system.listGroups", []interface{}{sid}, &result)
	if err != nil {
		return nil, fmt.Errorf("list groups: %w", err)
	}
	return result, nil
}

func (s *SystemService) ListActiveSystems(ctx context.Context) ([]models.System, error) {
	var result []models.System
	err := s.client(ctx, "system.listActiveSystems", nil, &result)
	if err != nil {
		return nil, fmt.Errorf("list active systems: %w", err)
	}
	return result, nil
}

func (s *SystemService) ListInactiveSystems(ctx context.Context, days int) ([]models.System, error) {
	var result []models.System
	if days > 0 {
		err := s.client(ctx, "system.listInactiveSystems", []interface{}{days}, &result)
		if err != nil {
			return nil, fmt.Errorf("list inactive systems: %w", err)
		}
	} else {
		err := s.client(ctx, "system.listInactiveSystems", nil, &result)
		if err != nil {
			return nil, fmt.Errorf("list inactive systems: %w", err)
		}
	}
	return result, nil
}

func (s *SystemService) ListOutOfDateSystems(ctx context.Context) ([]models.System, error) {
	var result []models.System
	err := s.client(ctx, "system.listOutOfDateSystems", nil, &result)
	if err != nil {
		return nil, fmt.Errorf("list out-of-date systems: %w", err)
	}
	return result, nil
}

func (s *SystemService) ListSuggestedReboot(ctx context.Context) ([]models.System, error) {
	var result []models.System
	err := s.client(ctx, "system.listSuggestedReboot", nil, &result)
	if err != nil {
		return nil, fmt.Errorf("list suggested reboot: %w", err)
	}
	return result, nil
}
