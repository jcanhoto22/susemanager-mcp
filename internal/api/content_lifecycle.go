package api

import (
	"context"
	"fmt"
)

type ContentLifecycleService struct {
	client func(ctx context.Context, method string, args []interface{}, reply interface{}) error
}

func NewContentLifecycleService(apiClient *Client) *ContentLifecycleService {
	return &ContentLifecycleService{
		client: apiClient.Call,
	}
}

func (s *ContentLifecycleService) ListProjects(ctx context.Context) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	err := s.client(ctx, "contentmanagement.listProjects", nil, &result)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	return result, nil
}

func (s *ContentLifecycleService) LookupProject(ctx context.Context, label string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := s.client(ctx, "contentmanagement.lookupProject", []interface{}{label}, &result)
	if err != nil {
		return nil, fmt.Errorf("lookup project: %w", err)
	}
	return result, nil
}

func (s *ContentLifecycleService) CreateProject(ctx context.Context, label, name, description string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := s.client(ctx, "contentmanagement.createProject", []interface{}{label, name, description}, &result)
	if err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}
	return result, nil
}

func (s *ContentLifecycleService) RemoveProject(ctx context.Context, label string) (int, error) {
	var result int
	err := s.client(ctx, "contentmanagement.removeProject", []interface{}{label}, &result)
	if err != nil {
		return 0, fmt.Errorf("remove project: %w", err)
	}
	return result, nil
}

func (s *ContentLifecycleService) BuildProject(ctx context.Context, label string, message string) (int, error) {
	var result int
	args := []interface{}{label}
	if message != "" {
		args = append(args, message)
	}
	err := s.client(ctx, "contentmanagement.buildProject", args, &result)
	if err != nil {
		return 0, fmt.Errorf("build project: %w", err)
	}
	return result, nil
}

func (s *ContentLifecycleService) PromoteProject(ctx context.Context, projectLabel, envLabel string) (int, error) {
	var result int
	err := s.client(ctx, "contentmanagement.promoteProject", []interface{}{projectLabel, envLabel}, &result)
	if err != nil {
		return 0, fmt.Errorf("promote project: %w", err)
	}
	return result, nil
}

func (s *ContentLifecycleService) ListProjectEnvironments(ctx context.Context, projectLabel string) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	err := s.client(ctx, "contentmanagement.listProjectEnvironments", []interface{}{projectLabel}, &result)
	if err != nil {
		return nil, fmt.Errorf("list project environments: %w", err)
	}
	return result, nil
}

func (s *ContentLifecycleService) CreateEnvironment(ctx context.Context, projectLabel, predecessorLabel, envLabel, name, description string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := s.client(ctx, "contentmanagement.createEnvironment", []interface{}{projectLabel, predecessorLabel, envLabel, name, description}, &result)
	if err != nil {
		return nil, fmt.Errorf("create environment: %w", err)
	}
	return result, nil
}

func (s *ContentLifecycleService) AttachSource(ctx context.Context, projectLabel, sourceType, sourceLabel string, position int) (map[string]interface{}, error) {
	var result map[string]interface{}
	args := []interface{}{projectLabel, sourceType, sourceLabel}
	if position > 0 {
		args = append(args, position)
	}
	err := s.client(ctx, "contentmanagement.attachSource", args, &result)
	if err != nil {
		return nil, fmt.Errorf("attach source: %w", err)
	}
	return result, nil
}

func (s *ContentLifecycleService) DetachSource(ctx context.Context, projectLabel, sourceType, sourceLabel string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := s.client(ctx, "contentmanagement.detachSource", []interface{}{projectLabel, sourceType, sourceLabel}, &result)
	if err != nil {
		return nil, fmt.Errorf("detach source: %w", err)
	}
	return result, nil
}

func (s *ContentLifecycleService) ListProjectSources(ctx context.Context, projectLabel, sourceType string) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	args := []interface{}{projectLabel}
	if sourceType != "" {
		args = append(args, sourceType)
	}
	err := s.client(ctx, "contentmanagement.listProjectSources", args, &result)
	if err != nil {
		return nil, fmt.Errorf("list project sources: %w", err)
	}
	return result, nil
}

func (s *ContentLifecycleService) SetSources(ctx context.Context, projectLabel, sourceType string, sourceLabels []string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := s.client(ctx, "contentmanagement.setSources", []interface{}{projectLabel, sourceType, sourceLabels}, &result)
	if err != nil {
		return nil, fmt.Errorf("set sources: %w", err)
	}
	return result, nil
}
