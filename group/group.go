package group

import (
	"context"
	"fmt"

	group "google.golang.org/api/groupssettings/v1"
	"google.golang.org/api/option"
)

type GroupService struct {
	service *group.Service
}

func NewGroupService(ctx context.Context) (*GroupService, error) {
	srv, err := group.NewService(ctx, option.WithCredentialsFile("private.json"), option.WithScopes(
		group.AppsGroupsSettingsScope,
	))
	if err != nil {
		return nil, fmt.Errorf("failed to start group service, error: %v", err)
	}
	return &GroupService{
		service: srv,
	}, nil
}

type CreateGroupRequest struct {
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

func (gs *GroupService) CreateGroup(ctx context.Context, req *CreateGroupRequest) (string, error) {
	gs.service.Groups.Update("", &group.Groups{})
	return "", nil
}
