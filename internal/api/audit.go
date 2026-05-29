package api

import (
	"context"
	"fmt"

	"mcp-susemanager/internal/models"
)

type AuditService struct {
	client func(ctx context.Context, method string, args []interface{}, reply interface{}) error
}

func NewAuditService(apiClient *Client) *AuditService {
	return &AuditService{
		client: apiClient.Call,
	}
}

func (s *AuditService) ListSystemsByPatchStatus(ctx context.Context, cveIdentifier string, patchStatuses []string) ([]models.CveAuditSystem, error) {
	var result []models.CveAuditSystem
	args := []interface{}{cveIdentifier}
	if len(patchStatuses) > 0 {
		args = append(args, patchStatuses)
		err := s.client(ctx, "audit.listSystemsByPatchStatus", args, &result)
		if err != nil {
			return nil, fmt.Errorf("list systems by patch status: %w", err)
		}
	} else {
		err := s.client(ctx, "audit.listSystemsByPatchStatus", args, &result)
		if err != nil {
			return nil, fmt.Errorf("list systems by patch status: %w", err)
		}
	}
	return result, nil
}
