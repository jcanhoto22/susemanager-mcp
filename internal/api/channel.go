package api

import (
	"context"
	"fmt"

	"mcp-susemanager/internal/models"
)

type ChannelService struct {
	client func(ctx context.Context, method string, args []interface{}, reply interface{}) error
}

func NewChannelService(apiClient *Client) *ChannelService {
	return &ChannelService{
		client: apiClient.Call,
	}
}

func (s *ChannelService) ListSoftwareChannels(ctx context.Context) ([]models.Channel, error) {
	var result []models.Channel
	err := s.client(ctx, "channel.listSoftwareChannels", nil, &result)
	if err != nil {
		return nil, fmt.Errorf("list software channels: %w", err)
	}
	return result, nil
}

func (s *ChannelService) ListAllChannels(ctx context.Context) ([]models.ChannelListItem, error) {
	var result []models.ChannelListItem
	err := s.client(ctx, "channel.listAllChannels", nil, &result)
	if err != nil {
		return nil, fmt.Errorf("list all channels: %w", err)
	}
	return result, nil
}

func (s *ChannelService) GetDetails(ctx context.Context, label string) (*models.ChannelDetails, error) {
	var result models.ChannelDetails
	err := s.client(ctx, "channel.software.getDetails", []interface{}{label}, &result)
	if err != nil {
		return nil, fmt.Errorf("get channel details: %w", err)
	}
	return &result, nil
}

func (s *ChannelService) ListPackages(ctx context.Context, label string) ([]models.ChannelPackage, error) {
	var result []models.ChannelPackage
	err := s.client(ctx, "channel.software.listPackages", []interface{}{label}, &result)
	if err != nil {
		return nil, fmt.Errorf("list channel packages: %w", err)
	}
	return result, nil
}

func (s *ChannelService) ListSubscribedChildChannels(ctx context.Context, sid int) ([]models.Channel, error) {
	var result []models.Channel
	err := s.client(ctx, "system.listSubscribedChildChannels", []interface{}{sid}, &result)
	if err != nil {
		return nil, fmt.Errorf("list subscribed child channels: %w", err)
	}
	return result, nil
}

func (s *ChannelService) ListSubscribableBaseChannels(ctx context.Context, sid int) ([]models.Channel, error) {
	var result []models.Channel
	err := s.client(ctx, "system.listSubscribableBaseChannels", []interface{}{sid}, &result)
	if err != nil {
		return nil, fmt.Errorf("list subscribable base channels: %w", err)
	}
	return result, nil
}

func (s *ChannelService) SetChildChannels(ctx context.Context, sid int, channelIDs []int) (int, error) {
	var result int
	err := s.client(ctx, "system.setChildChannels", []interface{}{sid, channelIDs}, &result)
	if err != nil {
		return 0, fmt.Errorf("set child channels: %w", err)
	}
	return result, nil
}

func (s *ChannelService) SetBaseChannel(ctx context.Context, sid int, channelLabel string) (int, error) {
	var result int
	err := s.client(ctx, "system.setBaseChannel", []interface{}{sid, channelLabel}, &result)
	if err != nil {
		return 0, fmt.Errorf("set base channel: %w", err)
	}
	return result, nil
}

func (s *ChannelService) CreateRepo(ctx context.Context, label, repoType, url string) (*models.RepoInfo, error) {
	var result models.RepoInfo
	err := s.client(ctx, "channel.software.createRepo", []interface{}{label, repoType, url}, &result)
	if err != nil {
		return nil, fmt.Errorf("create repo: %w", err)
	}
	return &result, nil
}

func (s *ChannelService) Create(ctx context.Context, label, name, summary, archLabel, parentLabel string) (int, error) {
	var result int
	err := s.client(ctx, "channel.software.create", []interface{}{label, name, summary, archLabel, parentLabel}, &result)
	if err != nil {
		return 0, fmt.Errorf("create channel: %w", err)
	}
	return result, nil
}

func (s *ChannelService) AssociateRepo(ctx context.Context, channelLabel, repoLabel string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := s.client(ctx, "channel.software.associateRepo", []interface{}{channelLabel, repoLabel}, &result)
	if err != nil {
		return nil, fmt.Errorf("associate repo: %w", err)
	}
	return result, nil
}

func (s *ChannelService) ListArches(ctx context.Context) ([]map[string]string, error) {
	var result []map[string]string
	err := s.client(ctx, "channel.software.listArches", nil, &result)
	if err != nil {
		return nil, fmt.Errorf("list arches: %w", err)
	}
	return result, nil
}

func (s *ChannelService) DeleteChannel(ctx context.Context, channelLabel string) (int, error) {
	var result int
	err := s.client(ctx, "channel.software.delete", []interface{}{channelLabel}, &result)
	if err != nil {
		return 0, fmt.Errorf("delete channel: %w", err)
	}
	return result, nil
}

func (s *ChannelService) RemoveRepo(ctx context.Context, label string) (int, error) {
	var result int
	err := s.client(ctx, "channel.software.removeRepo", []interface{}{label}, &result)
	if err != nil {
		return 0, fmt.Errorf("remove repo: %w", err)
	}
	return result, nil
}

func (s *ChannelService) ListUserRepos(ctx context.Context) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	err := s.client(ctx, "channel.software.listUserRepos", nil, &result)
	if err != nil {
		return nil, fmt.Errorf("list user repos: %w", err)
	}
	return result, nil
}

func (s *ChannelService) ListChannelRepos(ctx context.Context, channelLabel string) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	err := s.client(ctx, "channel.software.listChannelRepos", []interface{}{channelLabel}, &result)
	if err != nil {
		return nil, fmt.Errorf("list channel repos: %w", err)
	}
	return result, nil
}

func (s *ChannelService) ScheduleChangeChannels(ctx context.Context, sids []int, baseChannelLabel string, childLabels []string, earliest string) (int, error) {
	var result int
	args := []interface{}{sids, baseChannelLabel, childLabels}
	if earliest != "" {
		args = append(args, earliest)
	}
	err := s.client(ctx, "system.scheduleChangeChannels", args, &result)
	if err != nil {
		return 0, fmt.Errorf("schedule change channels: %w", err)
	}
	return result, nil
}
