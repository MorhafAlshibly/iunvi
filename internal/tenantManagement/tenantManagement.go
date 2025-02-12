package tenantManagement

import (
	"context"

	"connectrpc.com/connect"
	"github.com/MorhafAlshibly/iunvi/gen/api"
	"github.com/MorhafAlshibly/iunvi/internal/tenantManagement/model"
)

type Service struct {
	database *model.Queries
}

func WithDatabase(database *model.Queries) func(*Service) {
	return func(input *Service) {
		input.database = database
	}
}

func NewService(options ...func(*Service)) *Service {
	service := &Service{}
	for _, option := range options {
		option(service)
	}
	return service
}

func (s *Service) CreateWorkspace(ctx context.Context, req *connect.Request[api.CreateWorkspaceRequest]) (*connect.Response[api.CreateWorkspaceResponse], error) {
	_, err := s.database.CreateWorkspace(ctx, req.Msg.Name)
	if err != nil {
		return nil, err
	}
	workspace, err := s.database.GetWorkspace(ctx, req.Msg.Name)
	if err != nil {
		return nil, err
	}
	res := connect.NewResponse(&api.CreateWorkspaceResponse{
		Id: workspace.WorkspaceID.String(),
	})
	return res, nil
}

func unmarshalWorkspace(workspace *model.AuthWorkspace) *api.Workspace {
	return &api.Workspace{
		Id:   workspace.WorkspaceID.String(),
		Name: workspace.Name,
	}
}

func unmarshalWorkspaces(workspaces []model.AuthWorkspace) []*api.Workspace {
	res := make([]*api.Workspace, len(workspaces))
	for i, workspace := range workspaces {
		res[i] = unmarshalWorkspace(&workspace)
	}
	return res
}

func (s *Service) GetWorkspaces(ctx context.Context, req *connect.Request[api.GetWorkspacesRequest]) (*connect.Response[api.GetWorkspacesResponse], error) {
	workspaces, err := s.database.GetWorkspaces(ctx, model.GetWorkspacesParams{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		return nil, err
	}
	res := connect.NewResponse(&api.GetWorkspacesResponse{
		Workspaces: unmarshalWorkspaces(workspaces),
	})
	return res, nil
}
