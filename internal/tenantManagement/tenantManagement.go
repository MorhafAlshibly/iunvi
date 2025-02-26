package tenantManagement

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/containers/azcontainerregistry"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry"
	"github.com/MorhafAlshibly/iunvi/gen/api"
	"github.com/MorhafAlshibly/iunvi/internal/tenantManagement/model"
	"github.com/MorhafAlshibly/iunvi/pkg/authorization"
	"github.com/MorhafAlshibly/iunvi/pkg/conversion"
	"github.com/MorhafAlshibly/iunvi/pkg/middleware"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	subscriptionId    string
	resourceGroupName string
	tenantId          string
	clientId          string
	clientSecret      string
	registryName      string
	registryTokenName string
}

func WithSubscriptionId(subscriptionId string) func(*Service) {
	return func(input *Service) {
		input.subscriptionId = subscriptionId
	}
}

func WithResourceGroupName(resourceGroupName string) func(*Service) {
	return func(input *Service) {
		input.resourceGroupName = resourceGroupName
	}
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

func WithRegistryName(registryName string) func(*Service) {
	return func(input *Service) {
		input.registryName = registryName
	}
}

func WithRegistryTokenName(registryTokenName string) func(*Service) {
	return func(input *Service) {
		input.registryTokenName = registryTokenName
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
	credentials, err := azidentity.NewClientSecretCredential(s.tenantId, s.clientId, s.clientSecret, nil)
	if err != nil {
		return nil, err
	}
	scopeMapClient, err := armcontainerregistry.NewScopeMapsClient(s.subscriptionId, credentials, nil)
	if err != nil {
		return nil, err
	}
	scope := getScope(workspace.WorkspaceID.String())
	poller, err := scopeMapClient.BeginCreate(ctx, s.resourceGroupName, s.registryName, scope, armcontainerregistry.ScopeMap{
		Properties: &armcontainerregistry.ScopeMapProperties{
			// Only read, list and write actions are allowed
			Actions: []*string{
				conversion.ValueToPointer(fmt.Sprintf("repositories/%s/*/%s", scope, "content/write")),
				conversion.ValueToPointer(fmt.Sprintf("repositories/%s/*/%s", scope, "metadata/read")),
				conversion.ValueToPointer(fmt.Sprintf("repositories/%s/*/%s", scope, "metadata/write")),
			},
			// Assign the scope to the workspace
			Description: conversion.ValueToPointer(fmt.Sprintf("Scope for workspace %s", scope)),
		},
	}, nil)
	if err != nil {
		return nil, err
	}
	scopeMapResult, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	tokenClient, err := armcontainerregistry.NewTokensClient(s.subscriptionId, credentials, nil)
	if err != nil {
		return nil, err
	}
	tokenName := s.getTokenName(workspace.WorkspaceID.String())
	_, err = tokenClient.BeginCreate(ctx, s.resourceGroupName, s.registryName, tokenName, armcontainerregistry.Token{
		Properties: &armcontainerregistry.TokenProperties{
			Credentials: nil,
			ScopeMapID:  scopeMapResult.ID,
			Status:      conversion.ValueToPointer(armcontainerregistry.TokenStatusEnabled),
		},
	}, nil)
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
	if req.Msg.Id == "" {
		return nil, fmt.Errorf("id is required")
	}
	database := model.New(middleware.GetTx(ctx))
	workspaceIdBytes, err := conversion.StringToUniqueIdentifier(req.Msg.Id)
	if err != nil {
		return nil, err
	}
	_, err = database.EditWorkspace(ctx, model.EditWorkspaceParams{
		WorkspaceId: workspaceIdBytes,
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
	if err != nil {
		return nil, err
	}
	var users []*api.User
	err = pageIterator.Iterate(context.Background(), func(user models.Userable) bool {
		id := conversion.PointerToValue(user.GetId(), "00000000-0000-0000-0000-000000000000")
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

func (s *Service) GetUserWorkspaceAssignment(ctx context.Context, req *connect.Request[api.GetUserWorkspaceAssignmentRequest]) (*connect.Response[api.GetUserWorkspaceAssignmentResponse], error) {
	if req.Msg.UserObjectId == "" {
		return nil, fmt.Errorf("UserObjectId is required")
	}
	if req.Msg.WorkspaceId == "" {
		return nil, fmt.Errorf("WorkspaceId is required")
	}
	userObjectIdBytes, err := conversion.StringToUniqueIdentifier(req.Msg.UserObjectId)
	if err != nil {
		return nil, err
	}
	workspaceIdBytes, err := conversion.StringToUniqueIdentifier(req.Msg.WorkspaceId)
	if err != nil {
		return nil, err
	}
	database := model.New(middleware.GetTx(ctx))
	assignment, err := database.GetUserWorkspaceAssignment(ctx, model.GetUserWorkspaceAssignmentParams{
		UserObjectId: userObjectIdBytes,
		WorkspaceId:  workspaceIdBytes,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return connect.NewResponse(&api.GetUserWorkspaceAssignmentResponse{
				Role: api.WorkspaceRole_UNASSIGNED,
			}), nil
		}
		return nil, err
	}
	if assignment.RoleName == "Viewer" {
		return connect.NewResponse(&api.GetUserWorkspaceAssignmentResponse{
			Role: api.WorkspaceRole_VIEWER,
		}), nil
	}
	if assignment.RoleName == "User" {
		return connect.NewResponse(&api.GetUserWorkspaceAssignmentResponse{
			Role: api.WorkspaceRole_USER,
		}), nil
	}
	if assignment.RoleName == "Developer" {
		return connect.NewResponse(&api.GetUserWorkspaceAssignmentResponse{
			Role: api.WorkspaceRole_DEVELOPER,
		}), nil
	}
	return nil, fmt.Errorf("unknown role: %s", assignment.RoleName)
}

func (s *Service) AssignUserToWorkspace(ctx context.Context, req *connect.Request[api.AssignUserToWorkspaceRequest]) (*connect.Response[api.AssignUserToWorkspaceResponse], error) {
	if req.Msg.UserObjectId == "" {
		return nil, fmt.Errorf("UserObjectId is required")
	}
	if req.Msg.WorkspaceId == "" {
		return nil, fmt.Errorf("WorkspaceId is required")
	}
	userObjectIdBytes, err := conversion.StringToUniqueIdentifier(req.Msg.UserObjectId)
	if err != nil {
		return nil, err
	}
	workspaceIdBytes, err := conversion.StringToUniqueIdentifier(req.Msg.WorkspaceId)
	if err != nil {
		return nil, err
	}
	database := model.New(middleware.GetTx(ctx))
	_, err = database.DeleteUserWorkspaceAssignment(ctx, model.DeleteUserWorkspaceAssignmentParams{
		UserObjectId: userObjectIdBytes,
		WorkspaceId:  workspaceIdBytes,
	})
	if err != nil {
		return nil, err
	}
	if req.Msg.Role != api.WorkspaceRole_UNASSIGNED {
		_, err = database.AssignUserToWorkspace(ctx, model.AssignUserToWorkspaceParams{
			UserObjectId: userObjectIdBytes,
			WorkspaceId:  workspaceIdBytes,
			RoleName:     req.Msg.Role.String(),
		})
		if err != nil {
			return nil, err
		}
	}
	return connect.NewResponse(&api.AssignUserToWorkspaceResponse{}), nil
}

func (s *Service) GetRegistryTokenPasswords(ctx context.Context, req *connect.Request[api.GetRegistryTokenPasswordsRequest]) (*connect.Response[api.GetRegistryTokenPasswordsResponse], error) {
	if req.Msg.WorkspaceId == "" {
		return nil, fmt.Errorf("WorkspaceId is required")
	}
	workspaceIdBytes, err := conversion.StringToUniqueIdentifier(req.Msg.WorkspaceId)
	if err != nil {
		return nil, err
	}
	authorization, err := authorization.NewAuthorization(authorization.WithWorkspaceID(workspaceIdBytes), authorization.WithWorkspaceRole(api.WorkspaceRole_DEVELOPER)).IsAuthorized(ctx)
	if err != nil {
		return nil, err
	}
	if !authorization {
		return nil, fmt.Errorf("unauthorized")
	}
	credentials, err := azidentity.NewClientSecretCredential(s.tenantId, s.clientId, s.clientSecret, nil)
	if err != nil {
		return nil, err
	}
	tokenClient, err := armcontainerregistry.NewTokensClient(s.subscriptionId, credentials, nil)
	if err != nil {
		return nil, err
	}
	tokenName := s.getTokenName(req.Msg.WorkspaceId)
	token, err := tokenClient.Get(ctx, s.resourceGroupName, s.registryName, tokenName, nil)
	if err != nil {
		return nil, err
	}
	var password1 *timestamppb.Timestamp
	var password2 *timestamppb.Timestamp
	if token.Properties.Credentials != nil {
		if len(token.Properties.Credentials.Passwords) >= 1 {
			password1Time := token.Properties.Credentials.Passwords[0].CreationTime
			if password1Time != nil {
				password1 = timestamppb.New(*password1Time)
			}
		}
		if len(token.Properties.Credentials.Passwords) >= 2 {
			password2Time := token.Properties.Credentials.Passwords[1].CreationTime
			if password2Time != nil {
				password2 = timestamppb.New(*password2Time)
			}
		}
	}
	return connect.NewResponse(&api.GetRegistryTokenPasswordsResponse{
		Password1: password1,
		Password2: password2,
	}), nil
}

func (s *Service) CreateRegistryTokenPassword(ctx context.Context, req *connect.Request[api.CreateRegistryTokenPasswordRequest]) (*connect.Response[api.CreateRegistryTokenPasswordResponse], error) {
	if req.Msg.WorkspaceId == "" {
		return nil, fmt.Errorf("WorkspaceId is required")
	}
	workspaceIdBytes, err := conversion.StringToUniqueIdentifier(req.Msg.WorkspaceId)
	if err != nil {
		return nil, err
	}
	authorization, err := authorization.NewAuthorization(authorization.WithWorkspaceID(workspaceIdBytes), authorization.WithWorkspaceRole(api.WorkspaceRole_DEVELOPER)).IsAuthorized(ctx)
	if err != nil {
		return nil, err
	}
	if !authorization {
		return nil, fmt.Errorf("unauthorized")
	}
	credentials, err := azidentity.NewClientSecretCredential(s.tenantId, s.clientId, s.clientSecret, nil)
	if err != nil {
		return nil, err
	}
	tokenClient, err := armcontainerregistry.NewTokensClient(s.subscriptionId, credentials, nil)
	if err != nil {
		return nil, err
	}
	tokenName := s.getTokenName(req.Msg.WorkspaceId)
	passwordName := armcontainerregistry.TokenPasswordNamePassword1
	passwordIndex := 0
	if req.Msg.Password2 {
		passwordName = armcontainerregistry.TokenPasswordNamePassword2
		passwordIndex = 1
	}
	token, err := tokenClient.Get(ctx, s.resourceGroupName, s.registryName, tokenName, nil)
	if err != nil {
		return nil, err
	}
	clientFactory, err := armcontainerregistry.NewClientFactory(s.subscriptionId, credentials, nil)
	if err != nil {
		return nil, err
	}
	poller, err := clientFactory.NewRegistriesClient().BeginGenerateCredentials(ctx, s.resourceGroupName, s.registryName, armcontainerregistry.GenerateCredentialsParameters{
		TokenID: token.ID,
		Name:    &passwordName,
	}, nil)
	if err != nil {
		return nil, err
	}
	result, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	if len(result.Passwords) <= passwordIndex {
		return nil, fmt.Errorf("password not returned")
	}
	password := result.Passwords[passwordIndex].Value
	if password == nil {
		return nil, fmt.Errorf("password is nil")
	}
	createdAt := result.Passwords[passwordIndex].CreationTime
	if createdAt == nil {
		return nil, fmt.Errorf("creation time is nil")
	}
	createdAtTimestamp := timestamppb.New(*createdAt)
	return connect.NewResponse(&api.CreateRegistryTokenPasswordResponse{
		Password:  *password,
		CreatedAt: createdAtTimestamp,
	}), nil
}

func (s *Service) GetImages(ctx context.Context, req *connect.Request[api.GetImagesRequest]) (*connect.Response[api.GetImagesResponse], error) {
	if req.Msg.WorkspaceId == "" {
		return nil, fmt.Errorf("WorkspaceId is required")
	}
	workspaceIdBytes, err := conversion.StringToUniqueIdentifier(req.Msg.WorkspaceId)
	if err != nil {
		return nil, err
	}
	authorization, err := authorization.NewAuthorization(authorization.WithWorkspaceID(workspaceIdBytes), authorization.WithWorkspaceRole(api.WorkspaceRole_VIEWER)).IsAuthorized(ctx)
	if err != nil {
		return nil, err
	}
	if !authorization {
		return nil, fmt.Errorf("unauthorized")
	}
	credentials, err := azidentity.NewClientSecretCredential(s.tenantId, s.clientId, s.clientSecret, nil)
	if err != nil {
		return nil, err
	}
	client, err := azcontainerregistry.NewClient(fmt.Sprintf("https://%s.azurecr.io", s.registryName), credentials, nil)
	if err != nil {
		return nil, err
	}
	scope := getScope(req.Msg.WorkspaceId)
	pager := client.NewListRepositoriesPager(nil)
	var images []*api.Image
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, repository := range page.Repositories.Names {
			if repository == nil {
				continue
			}
			if !strings.HasPrefix(*repository, scope) {
				continue
			}
			images = append(images, &api.Image{
				Name: *repository,
			})
		}
	}
	return connect.NewResponse(&api.GetImagesResponse{
		Images: images,
	}), nil
}

func getScope(workspaceId string) string {
	return fmt.Sprintf("scope-%s", strings.ToLower(workspaceId))
}

func (s *Service) getTokenName(workspaceId string) string {
	return fmt.Sprintf("%s-%s", s.registryTokenName, strings.ToLower(workspaceId))
}
