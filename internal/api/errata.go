package api

import (
	"context"
	"fmt"

	"mcp-susemanager/internal/models"
)

type ErrataService struct {
	client func(ctx context.Context, method string, args []interface{}, reply interface{}) error
}

func NewErrataService(apiClient *Client) *ErrataService {
	return &ErrataService{
		client: apiClient.Call,
	}
}

func (s *ErrataService) GetRelevantErrata(ctx context.Context, sid int) ([]models.Errata, error) {
	var result []models.Errata
	err := s.client(ctx, "system.getRelevantErrata", []interface{}{sid}, &result)
	if err != nil {
		return nil, fmt.Errorf("get relevant errata: %w", err)
	}
	return result, nil
}

func (s *ErrataService) GetUnscheduledErrata(ctx context.Context, sid int) ([]models.Errata, error) {
	var result []models.Errata
	err := s.client(ctx, "system.getUnscheduledErrata", []interface{}{sid}, &result)
	if err != nil {
		return nil, fmt.Errorf("get unscheduled errata: %w", err)
	}
	return result, nil
}

func (s *ErrataService) ScheduleApplyErrata(ctx context.Context, sids []int, errataIDs []int, earliest string) (int, error) {
	var result int
	args := []interface{}{sids, errataIDs}
	if earliest != "" {
		args = append(args, earliest)
	}
	err := s.client(ctx, "system.scheduleApplyErrata", args, &result)
	if err != nil {
		return 0, fmt.Errorf("schedule apply errata: %w", err)
	}
	return result, nil
}
