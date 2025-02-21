package tenantManagement

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/MorhafAlshibly/iunvi/gen/api"
	"github.com/MorhafAlshibly/iunvi/internal/tenantManagement/model"
	"github.com/MorhafAlshibly/iunvi/pkg/conversion"
	"github.com/MorhafAlshibly/iunvi/pkg/middleware"
	"github.com/google/uuid"
	mssql "github.com/microsoft/go-mssqldb"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

type Service struct {
	tenantId     string
	clientId     string
	clientSecret string
}

func WithTenantId(tenantId string) func(*Service) {
	return func(input *Service) {
		input.tenantId = tenantId
	}
}

func WithClientId(clientId string) func(*Service) {
	return func(input *Service) {
		input.clientId = clientId
	}
}

func WithClientSecret(clientSecret string) func(*Service) {
	return func(input *Service) {
		input.clientSecret = clientSecret
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
	database := model.New(middleware.GetTx(ctx))
	_, err := database.CreateWorkspace(ctx, req.Msg.Name)
	if err != nil {
		return nil, err
	}
	workspace, err := database.GetWorkspace(ctx, req.Msg.Name)
	if err != nil {
		return nil, err
	}
	res := connect.NewResponse(&api.CreateWorkspaceResponse{
		Id: workspace.WorkspaceID[:],
	})
	return res, nil
}

func unmarshalWorkspace(workspace *model.AuthWorkspace) *api.Workspace {
	return &api.Workspace{
		Id:   workspace.WorkspaceID[:],
		Name: workspace.Name,
	}
}

func (s *Service) GetWorkspaces(ctx context.Context, req *connect.Request[api.GetWorkspacesRequest]) (*connect.Response[api.GetWorkspacesResponse], error) {
	database := model.New(middleware.GetTx(ctx))
	workspaces, err := database.GetWorkspaces(ctx, model.GetWorkspacesParams{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		return nil, err
	}
	marshalledWorkspaces := make([]*api.Workspace, len(workspaces))
	for i, workspace := range workspaces {
		marshalledWorkspaces[i] = unmarshalWorkspace(&workspace)
	}
	res := connect.NewResponse(&api.GetWorkspacesResponse{
		Workspaces: marshalledWorkspaces,
	})
	return res, nil
}

func (s *Service) EditWorkspace(ctx context.Context, req *connect.Request[api.EditWorkspaceRequest]) (*connect.Response[api.EditWorkspaceResponse], error) {
	database := model.New(middleware.GetTx(ctx))
	_, err := database.EditWorkspace(ctx, model.EditWorkspaceParams{
		WorkspaceId: mssql.UniqueIdentifier(req.Msg.Id),
		Name:        req.Msg.Name,
	})
	if err != nil {
		return nil, err
	}
	res := connect.NewResponse(&api.EditWorkspaceResponse{})
	return res, nil
}

func (s *Service) GetUsers(ctx context.Context, req *connect.Request[api.GetUsersRequest]) (*connect.Response[api.GetUsersResponse], error) {
	accessToken := ctx.Value("accessToken").(string)
	cred, err := azidentity.NewOnBehalfOfCredentialWithSecret(s.tenantId, s.clientId, accessToken, s.clientSecret, nil)
	if err != nil {
		return nil, err
	}
	client, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, []string{"User.ReadBasic.All"})
	if err != nil {
		return nil, err
	}
	result, err := client.Users().Get(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	pageIterator, err := msgraphcore.NewPageIterator[models.Userable](result, client.GetAdapter(), models.CreateUserCollectionResponseFromDiscriminatorValue)
	var users []*api.User
	err = pageIterator.Iterate(context.Background(), func(user models.Userable) bool {
		guid, err := uuid.Parse(conversion.PointerToValue(user.GetId(), "00000000-0000-0000-0000-000000000000"))
		if err != nil {
			return false
		}
		id, err := guid.MarshalBinary()
		if err != nil {
			return false
		}
		apiUser := &api.User{
			Id:          id,
			Username:    conversion.PointerToValue(user.GetUserPrincipalName(), ""),
			DisplayName: conversion.PointerToValue(user.GetDisplayName(), ""),
		}
		users = append(users, apiUser)
		return true
	})
	if err != nil {
		fmt.Printf("Error iterating through users: %v\n", err)
		return nil, err
	}
	return connect.NewResponse(&api.GetUsersResponse{
		Users: users,
	}), nil
}
