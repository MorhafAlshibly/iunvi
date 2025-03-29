package tenantManagement

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/containers/azcontainerregistry"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry"
	"github.com/MorhafAlshibly/iunvi/gen/api"
	"github.com/MorhafAlshibly/iunvi/internal/tenantManagement/model"
	"github.com/MorhafAlshibly/iunvi/pkg/authorization"
	"github.com/MorhafAlshibly/iunvi/pkg/conversion"
	"github.com/MorhafAlshibly/iunvi/pkg/middleware"
	mssql "github.com/microsoft/go-mssqldb"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"google.golang.org/protobuf/types/known/timestamppb"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/sas"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/service"
	_ "github.com/marcboeker/go-duckdb"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

type Service struct {
	subscriptionId                  string
	resourceGroupName               string
	tenantId                        string
	clientId                        string
	clientSecret                    string
	registryName                    string
	registryTokenName               string
	storageAccountName              string
	landingZoneContainerName        string
	fileGroupsContainerName         string
	modelRunsContainerName          string
	dashboardsContainerName         string
	modelRunDashboardsContainerName string
	kubeconfigPath                  string
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

func WithStorageAccountName(storageAccountName string) func(*Service) {
	return func(input *Service) {
		input.storageAccountName = storageAccountName
	}
}

func WithLandingZoneContainerName(landingZoneContainerName string) func(*Service) {
	return func(input *Service) {
		input.landingZoneContainerName = landingZoneContainerName
	}
}

func WithFileGroupsContainerName(fileGroupsContainerName string) func(*Service) {
	return func(input *Service) {
		input.fileGroupsContainerName = fileGroupsContainerName
	}
}

func WithModelRunsContainerName(modelRunsContainerName string) func(*Service) {
	return func(input *Service) {
		input.modelRunsContainerName = modelRunsContainerName
	}
}

func WithDashboardsContainerName(dashboardsContainerName string) func(*Service) {
	return func(input *Service) {
		input.dashboardsContainerName = dashboardsContainerName
	}
}

func WithModelRunDashboardsContainerName(modelRunDashboardsContainerName string) func(*Service) {
	return func(input *Service) {
		input.modelRunDashboardsContainerName = modelRunDashboardsContainerName
	}
}

func WithKubeconfigPath(kubeconfigPath string) func(*Service) {
	return func(input *Service) {
		input.kubeconfigPath = kubeconfigPath
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
				conversion.ValueToPointer(fmt.Sprintf("repositories/%s/*/%s", scope, "content/read")),
				conversion.ValueToPointer(fmt.Sprintf("repositories/%s/*/%s", scope, "content/write")),
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
			name := strings.TrimPrefix(*repository, scope+"/")
			images = append(images, &api.Image{
				Scope: scope,
				Name:  name,
			})
		}
	}
	return connect.NewResponse(&api.GetImagesResponse{
		Images: images,
	}), nil
}

func (s *Service) CreateSpecification(ctx context.Context, req *connect.Request[api.CreateSpecificationRequest]) (*connect.Response[api.CreateSpecificationResponse], error) {
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
	dataModeName := "Input"
	if req.Msg.Mode == api.DataMode_OUTPUT {
		dataModeName = "Output"
	}
	database := model.New(middleware.GetTx(ctx))
	_, err = database.CreateSpecification(ctx, model.CreateSpecificationParams{
		WorkspaceId:  workspaceIdBytes,
		DataModeName: dataModeName,
		Name:         req.Msg.Name,
	})
	if err != nil {
		return nil, err
	}
	specification, err := database.GetSpecificationByWorkspaceIdAndName(ctx, model.GetSpecificationByWorkspaceIdAndNameParams{
		WorkspaceId: workspaceIdBytes,
		Name:        req.Msg.Name,
	})
	if err != nil {
		return nil, err
	}
	for _, table := range req.Msg.Tables {
		if table.Name == "" {
			return nil, fmt.Errorf("name is required")
		}
		if table.Fields == nil || len(table.Fields) == 0 {
			return nil, fmt.Errorf("fields are required")
		}
		schema, err := json.Marshal(table.Fields)
		if err != nil {
			return nil, err
		}
		fileTypeName := "CSV"
		if req.Msg.Mode == api.DataMode_OUTPUT {
			fileTypeName = "Parquet"
		}
		_, err = database.CreateFileSchema(ctx, model.CreateFileSchemaParams{
			SpecificationId: specification.SpecificationID,
			FileTypeName:    fileTypeName,
			Name:            table.Name,
			Definition:      string(schema),
		})
		if err != nil {
			return nil, err
		}
	}
	return connect.NewResponse(&api.CreateSpecificationResponse{
		Id: specification.SpecificationID.String(),
	}), nil
}

func (s *Service) GetSpecifications(ctx context.Context, req *connect.Request[api.GetSpecificationsRequest]) (*connect.Response[api.GetSpecificationsResponse], error) {
	if req.Msg.WorkspaceId == "" {
		return nil, fmt.Errorf("WorkspaceId is required")
	}
	workspaceIdBytes, err := conversion.StringToUniqueIdentifier(req.Msg.WorkspaceId)
	if err != nil {
		return nil, err
	}
	database := model.New(middleware.GetTx(ctx))
	var mode *string
	if req.Msg.Mode != nil {
		mode = conversion.ValueToPointer(req.Msg.Mode.String())
	}
	specifications, err := database.GetSpecifications(ctx, model.GetSpecificationsParams{
		WorkspaceId:  workspaceIdBytes,
		DataModeName: mode,
	})
	if err != nil {
		return nil, err
	}
	var apiSpecifications []*api.SpecificationName
	for _, specification := range specifications {
		apiSpecifications = append(apiSpecifications, &api.SpecificationName{
			Id:   specification.SpecificationID.String(),
			Name: specification.Name,
		})
	}
	return connect.NewResponse(&api.GetSpecificationsResponse{
		Specifications: apiSpecifications,
	}), nil
}

func (s *Service) GetSpecification(ctx context.Context, req *connect.Request[api.GetSpecificationRequest]) (*connect.Response[api.GetSpecificationResponse], error) {
	if req.Msg.Id == "" {
		return nil, fmt.Errorf("SpecificationId is required")
	}
	specificationIdBytes, err := conversion.StringToUniqueIdentifier(req.Msg.Id)
	if err != nil {
		return nil, err
	}
	database := model.New(middleware.GetTx(ctx))
	specification, err := database.GetSpecification(ctx, model.GetSpecificationParams{
		SpecificationId: specificationIdBytes,
	})
	if err != nil {
		return nil, err
	}
	fileTypeName := "CSV"
	dataMode := api.DataMode_INPUT
	if specification.DataModeName == "Output" {
		fileTypeName = "Parquet"
		dataMode = api.DataMode_OUTPUT
	}
	var tables []*api.TableSchema
	tableSchemas, err := database.GetFileSchemasBySpecificationIdAndDataTypeName(ctx, model.GetFileSchemasBySpecificationIdAndDataTypeNameParams{
		SpecificationId: specificationIdBytes,
		FileTypeName:    fileTypeName,
	})
	if err != nil {
		return nil, err
	}
	if len(tableSchemas) == 0 {
		return nil, fmt.Errorf("table schemas not found")
	}
	for _, tableSchema := range tableSchemas {
		var fields []*api.TableField
		err := json.Unmarshal([]byte(tableSchema.Definition), &fields)
		if err != nil {
			return nil, err
		}
		table := &api.TableSchema{
			Name:   tableSchema.Name,
			Fields: fields,
		}
		tables = append(tables, table)
	}
	return connect.NewResponse(&api.GetSpecificationResponse{
		Specification: &api.Specification{
			Id:     specification.SpecificationID.String(),
			Name:   specification.Name,
			Tables: tables,
		},
		Mode: dataMode,
	}), nil
}

func (s *Service) CreateLandingZoneSharedAccessSignature(ctx context.Context, req *connect.Request[api.CreateLandingZoneSharedAccessSignatureRequest]) (*connect.Response[api.CreateLandingZoneSharedAccessSignatureResponse], error) {
	// if req.Msg.SpecificationID == "" {
	// 	return nil, fmt.Errorf("SpecificationId is required")
	// }
	// if req.Msg.Files == nil || len(req.Msg.Files) == 0 {
	// 	return nil, fmt.Errorf("Files are required")
	// }
	// // check for duplicate file names
	// fileNames := make(map[string]struct{})
	// fileSchemaIds := make(map[string]struct{})
	// for _, file := range req.Msg.Files {
	// 	fileSchemaIdBytes, err := conversion.StringToUniqueIdentifier(file.Id)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if _, ok := fileSchemaIds[file.Id]; ok {
	// 		return nil, fmt.Errorf("duplicate file schema id: %s", file.Id)
	// 	}
	// 	fileSchemaIds[file.Id] = struct{}{}

	// 	if _, ok := fileNames[file.Name]; ok {
	// 		return nil, fmt.Errorf("duplicate file name: %s", file.Name)
	// 	}
	// 	fileNames[file.Name] = struct{}{}
	// }
	// specificationIdBytes, err := conversion.StringToUniqueIdentifier(req.Msg.SpecificationID)
	// if err != nil {
	// 	return nil, err
	// }
	// database := model.New(middleware.GetTx(ctx))
	// tableSchemas, err := database.GetFileSchemasBySpecificationIdAndDataTypeName(ctx, model.GetFileSchemasBySpecificationIdAndDataTypeNameParams{
	// 	SpecificationId: specificationIdBytes,
	// 	FileTypeName:    "CSV",
	// })
	// if err != nil {
	// 	return nil, err
	// }
	// if len(tableSchemas) == 0 {
	// 	return nil, fmt.Errorf("table schemas not found")
	// }
	// // make sure every table schema has been requested
	// for _, tableSchema := range tableSchemas {
	// 	if _, ok := fileSchemaIds[tableSchema.FileSchemaID.String()]; !ok {
	// 		return nil, fmt.Errorf("missing table schema: %s", tableSchema.Name)
	// 	}
	// }
	// // make sure every requested file schema is a table schema
	// if len(tableSchemas) != len(fileSchemaIds) {
	// 	return nil, fmt.Errorf("invalid schema")
	// }
	// workspaceId, err := database.GetWorkspaceIdBySpecificationId(ctx, model.GetWorkspaceIdBySpecificationIdParams{
	// 	SpecificationId: specificationIdBytes,
	// })
	// if err != nil {
	// 	return nil, err
	// }
	if req.Msg.WorkspaceId == "" {
		return nil, fmt.Errorf("WorkspaceId is required")
	}
	if req.Msg.FileName == "" {
		return nil, fmt.Errorf("FileName is required")
	}
	fileExtension := strings.ToLower(req.Msg.FileName[strings.LastIndex(req.Msg.FileName, ".")+1:])
	if fileExtension != "csv" {
		return nil, fmt.Errorf("invalid file extension")
	}
	workspaceId, err := conversion.StringToUniqueIdentifier(req.Msg.WorkspaceId)
	if err != nil {
		return nil, err
	}
	authorization, err := authorization.NewAuthorization(authorization.WithWorkspaceID(workspaceId), authorization.WithWorkspaceRole(api.WorkspaceRole_USER)).IsAuthorized(ctx)
	if err != nil {
		return nil, err
	}
	if !authorization {
		return nil, fmt.Errorf("unauthorized")
	}
	cred, err := azidentity.NewClientSecretCredential(s.tenantId, s.clientId, s.clientSecret, nil)
	if err != nil {
		return nil, err
	}
	directory := getLandingZoneDirectory(ctx, req.Msg.WorkspaceId)
	svcClient, err := service.NewClient(fmt.Sprintf("https://%s.blob.core.windows.net", s.storageAccountName), cred, nil)
	if err != nil {
		return nil, err
	}
	info := service.KeyInfo{
		Start:  conversion.ValueToPointer(time.Now().UTC().Add(-10 * time.Second).Format(sas.TimeFormat)),
		Expiry: conversion.ValueToPointer(time.Now().UTC().Add(48 * time.Hour).Format(sas.TimeFormat)),
	}
	udc, err := svcClient.GetUserDelegationCredential(ctx, info, nil)
	if err != nil {
		return nil, err
	}
	// Create Blob Signature Values with desired permissions and sign with user delegation credential
	sasQueryParams, err := sas.BlobSignatureValues{
		Protocol:      sas.ProtocolHTTPS,
		StartTime:     time.Now().UTC().Add(-10 * time.Second),
		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
		Permissions:   (&sas.BlobPermissions{Write: true}).String(),
		ContainerName: s.landingZoneContainerName,
		BlobName:      directory + "/" + req.Msg.FileName,
	}.SignWithUserDelegation(udc)
	if err != nil {
		return nil, err
	}
	sasUrl := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s/%s?%s", s.storageAccountName, s.landingZoneContainerName, directory, req.Msg.FileName, sasQueryParams.Encode())
	return connect.NewResponse(&api.CreateLandingZoneSharedAccessSignatureResponse{
		Url: sasUrl,
	}), nil
}

func (s *Service) GetLandingZoneFiles(ctx context.Context, req *connect.Request[api.GetLandingZoneFilesRequest]) (*connect.Response[api.GetLandingZoneFilesResponse], error) {
	if req.Msg.WorkspaceId == "" {
		return nil, fmt.Errorf("WorkspaceId is required")
	}
	workspaceId, err := conversion.StringToUniqueIdentifier(req.Msg.WorkspaceId)
	if err != nil {
		return nil, err
	}
	authorization, err := authorization.NewAuthorization(authorization.WithWorkspaceID(workspaceId), authorization.WithWorkspaceRole(api.WorkspaceRole_USER)).IsAuthorized(ctx)
	if err != nil {
		return nil, err
	}
	if !authorization {
		return nil, fmt.Errorf("unauthorized")
	}
	cred, err := azidentity.NewClientSecretCredential(s.tenantId, s.clientId, s.clientSecret, nil)
	if err != nil {
		return nil, err
	}
	directory := getLandingZoneDirectory(ctx, req.Msg.WorkspaceId)
	svcClient, err := service.NewClient(fmt.Sprintf("https://%s.blob.core.windows.net", s.storageAccountName), cred, nil)
	if err != nil {
		return nil, err
	}
	containerClient := svcClient.NewContainerClient(s.landingZoneContainerName)
	pager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Prefix:     conversion.ValueToPointer(directory + "/" + req.Msg.Prefix),
		MaxResults: conversion.ValueToPointer(int32(10)),
		Marker:     req.Msg.Marker,
	})
	var files []*api.LandingZoneFile
	var page container.ListBlobsFlatResponse
	if pager.More() {
		page, err = pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, blob := range page.Segment.BlobItems {
			if blob.Name == nil {
				continue
			}
			fileName := strings.TrimPrefix(*blob.Name, directory+"/")
			files = append(files, &api.LandingZoneFile{
				Name:         fileName,
				Size:         uint64(conversion.PointerToValue(blob.Properties.ContentLength, 0)),
				LastModified: timestamppb.New(conversion.PointerToValue(blob.Properties.LastModified, time.Time{})),
			})
		}
	}
	if page.NextMarker != nil {
		if *page.NextMarker == "" {
			page.NextMarker = nil
		}
	}
	return connect.NewResponse(&api.GetLandingZoneFilesResponse{
		Files:      files,
		NextMarker: page.NextMarker,
	}), nil
}

func (s *Service) CreateFileGroup(ctx context.Context, req *connect.Request[api.CreateFileGroupRequest]) (*connect.Response[api.CreateFileGroupResponse], error) {
	if req.Msg.SpecificationId == "" {
		return nil, fmt.Errorf("SpecificationId is required")
	}
	if req.Msg.Name == "" {
		return nil, fmt.Errorf("Name is required")
	}
	if req.Msg.SchemaFileMappings == nil || len(req.Msg.SchemaFileMappings) == 0 {
		return nil, fmt.Errorf("SchemaFileMappings are required")
	}
	specificationIdBytes, err := conversion.StringToUniqueIdentifier(req.Msg.SpecificationId)
	if err != nil {
		return nil, err
	}
	database := model.New(middleware.GetTx(ctx))
	workspaceId, err := database.GetWorkspaceIdBySpecificationId(ctx, model.GetWorkspaceIdBySpecificationIdParams{
		SpecificationId: specificationIdBytes,
	})
	if err != nil {
		return nil, err
	}
	authorization, err := authorization.NewAuthorization(authorization.WithWorkspaceID(workspaceId), authorization.WithWorkspaceRole(api.WorkspaceRole_USER)).IsAuthorized(ctx)
	if err != nil {
		return nil, err
	}
	if !authorization {
		return nil, fmt.Errorf("unauthorized")
	}
	// make sure specification is for input data
	specification, err := database.GetSpecification(ctx, model.GetSpecificationParams{
		SpecificationId: specificationIdBytes,
	})
	if err != nil {
		return nil, err
	}
	if specification.DataModeName != "Input" {
		return nil, fmt.Errorf("invalid specification data mode")
	}
	tableSchemas, err := database.GetFileSchemasBySpecificationIdAndDataTypeName(ctx, model.GetFileSchemasBySpecificationIdAndDataTypeNameParams{
		SpecificationId: specificationIdBytes,
		FileTypeName:    "CSV",
	})
	if err != nil {
		return nil, err
	}
	if len(tableSchemas) == 0 {
		return nil, fmt.Errorf("table schemas not found")
	}
	userObjectId := ctx.Value("UserObjectId").(string)
	userObjectIdBytes, err := conversion.StringToUniqueIdentifier(userObjectId)
	if err != nil {
		return nil, err
	}
	_, err = database.CreateFileGroup(ctx, model.CreateFileGroupParams{
		SpecificationId:    specificationIdBytes,
		CreatedBy:          userObjectIdBytes,
		Name:               req.Msg.Name,
		ShareWithWorkspace: false,
	})
	if err != nil {
		return nil, err
	}
	fileGroup, err := database.GetFileGroupBySpecificationIdAndName(ctx, model.GetFileGroupBySpecificationIdAndNameParams{
		SpecificationId: specificationIdBytes,
		Name:            req.Msg.Name,
	})
	if err != nil {
		return nil, err
	}
	tableSchemasMap := make(map[string]model.GetFileSchemaBySpecificationIdAndDataTypeNameRow)
	for _, tableSchema := range tableSchemas {
		tableSchemasMap[tableSchema.Name] = tableSchema
	}
	cred, err := azidentity.NewClientSecretCredential(s.tenantId, s.clientId, s.clientSecret, nil)
	if err != nil {
		return nil, err
	}
	svcClient, err := service.NewClient(fmt.Sprintf("https://%s.blob.core.windows.net", s.storageAccountName), cred, nil)
	if err != nil {
		return nil, err
	}
	landingZoneContainerClient := svcClient.NewContainerClient(s.landingZoneContainerName)
	fileGroupsContainerClient := svcClient.NewContainerClient(s.fileGroupsContainerName)
	duckdb, err := sql.Open("duckdb", "")
	if err != nil {
		return nil, err
	}
	defer duckdb.Close()
	// loop through schema file mappings and check if the schema exists
	for _, schemaFileMapping := range req.Msg.SchemaFileMappings {
		if schemaFileMapping.SchemaName == "" {
			return nil, fmt.Errorf("SchemaName is required")
		}
		if schemaFileMapping.LandingZoneFileName == "" {
			return nil, fmt.Errorf("LandingZoneFileName is required")
		}
		if _, ok := tableSchemasMap[schemaFileMapping.SchemaName]; !ok {
			return nil, fmt.Errorf("schema not found: %s", schemaFileMapping.SchemaName)
		}
		var fields []*api.TableField
		err = json.Unmarshal([]byte(tableSchemasMap[schemaFileMapping.SchemaName].Definition), &fields)
		if err != nil {
			return nil, err
		}
		blobClient := landingZoneContainerClient.NewBlobClient(getLandingZoneDirectory(ctx, workspaceId.String()) + "/" + schemaFileMapping.LandingZoneFileName)
		// TODO: Acquire a lease on the blob to prevent other users from writing to it

		// Create []bytes of length 10KB to download the first 10KB of the file
		buffer := make([]byte, 10*1024)
		_, err = blobClient.DownloadBuffer(ctx, buffer, &blob.DownloadBufferOptions{
			Range: blob.HTTPRange{
				Offset: 0,
				Count:  10 * 1024,
			},
		})
		if err != nil {
			return nil, err
		}
		// Remove last bytes from buffer until we remove last newline character
		for i := len(buffer) - 1; i >= 0; i-- {
			if buffer[i] == '\n' {
				buffer = buffer[:i-1]
				break
			}
		}
		// Create a temporary file to store the first 10KB of the file
		tempFile, err := os.CreateTemp("", "temp")
		if err != nil {
			return nil, err
		}
		defer os.Remove(tempFile.Name())
		_, err = tempFile.Write(buffer)
		if err != nil {
			tempFile.Close()
			return nil, err
		}
		err = tempFile.Close()
		if err != nil {
			return nil, err
		}
		// Create a query to read the CSV file and get the schema
		columnsString := ""
		for i := 2; i < len(fields)+2; i++ {
			columnsString += "'$" + strconv.Itoa(i) + "': " + fields[i-2].Type.String() + ", "
		}
		query := "SELECT COUNT(*) FROM read_csv($1, max_line_size=10000000, columns={" + columnsString[:len(columnsString)-2] + "});"
		queryValues := make([]interface{}, 0, len(fields)+1)
		queryValues = append(queryValues, tempFile.Name())
		for _, field := range fields {
			queryValues = append(queryValues, field.Name)
		}
		rows, err := duckdb.Query(query, queryValues...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		// file is now validated
		// want to create a file in the file group
		_, err = database.CreateFile(ctx, model.CreateFileParams{
			FileGroupId:    fileGroup.FileGroupID,
			FileSchemaName: schemaFileMapping.SchemaName,
			Name:           schemaFileMapping.LandingZoneFileName,
		})
		if err != nil {
			return nil, err
		}
		// file, err := database.GetFileByFileGroupIdAndSchemaName(ctx, model.GetFileByFileGroupIdAndSchemaNameParams{
		// 	FileGroupId:    fileGroup.FileGroupID,
		// 	FileSchemaName: schemaFileMapping.LandingZoneFileName,
		// })
		// if err != nil {
		// 	return nil, err
		// }
	}
	// copy the files from the landing zone to the file group
	for _, schemaFileMapping := range req.Msg.SchemaFileMappings {
		landingZoneBlobClient := landingZoneContainerClient.NewBlobClient(getLandingZoneDirectory(ctx, workspaceId.String()) + "/" + schemaFileMapping.LandingZoneFileName)
		fileGroupsBlobClient := fileGroupsContainerClient.NewBlobClient(getFileGroupDirectory(ctx, workspaceId.String(), req.Msg.SpecificationId, fileGroup.FileGroupID.String()) + "/" + schemaFileMapping.LandingZoneFileName)
		_, err = fileGroupsBlobClient.StartCopyFromURL(ctx, landingZoneBlobClient.URL(), nil)
		if err != nil {
			return nil, err
		}
	}
	return connect.NewResponse(&api.CreateFileGroupResponse{
		Id: fileGroup.FileGroupID.String(),
	}), nil
}

func (s *Service) GetFileGroups(ctx context.Context, req *connect.Request[api.GetFileGroupsRequest]) (*connect.Response[api.GetFileGroupsResponse], error) {
	if req.Msg.WorkspaceId == "" {
		return nil, fmt.Errorf("WorkspaceId is required")
	}
	workspaceIdBytes, err := conversion.StringToUniqueIdentifier(req.Msg.WorkspaceId)
	if err != nil {
		return nil, err
	}
	var specificationIdBytes *mssql.UniqueIdentifier
	if req.Msg.SpecificationId != nil {
		specificationId, err := conversion.StringToUniqueIdentifier(*req.Msg.SpecificationId)
		if err != nil {
			return nil, err
		}
		specificationIdBytes = &specificationId
	}
	database := model.New(middleware.GetTx(ctx))
	fileGroups, err := database.GetFileGroups(ctx, model.GetFileGroupsParams{
		WorkspaceId:     workspaceIdBytes,
		SpecificationId: specificationIdBytes,
	})
	if err != nil {
		return nil, err
	}
	var fileGroupNames []*api.FileGroupName
	for _, fileGroup := range fileGroups {
		fileGroupNames = append(fileGroupNames, &api.FileGroupName{
			Id:   fileGroup.FileGroupID.String(),
			Name: fileGroup.Name,
		})
	}
	return connect.NewResponse(&api.GetFileGroupsResponse{
		FileGroups: fileGroupNames,
	}), nil
}

func (s *Service) CreateModel(ctx context.Context, req *connect.Request[api.CreateModelRequest]) (*connect.Response[api.CreateModelResponse], error) {
	if req.Msg.InputSpecificationId == "" {
		return nil, fmt.Errorf("InputSpecificationId is required")
	}
	if req.Msg.OutputSpecificationId == "" {
		return nil, fmt.Errorf("OutputSpecificationId is required")
	}
	if req.Msg.Name == "" {
		return nil, fmt.Errorf("Name is required")
	}
	var parametersSchema mssql.NVarCharMax
	if req.Msg.ParametersSchema != nil {
		if *req.Msg.ParametersSchema == "" {
			return nil, fmt.Errorf("ParametersSchema is required")
		}
		if len(*req.Msg.ParametersSchema) > 100000 {
			return nil, fmt.Errorf("ParametersSchema is too long")
		}
		// create temporary schema file
		tempFile, err := os.CreateTemp("", "temp")
		if err != nil {
			return nil, err
		}
		defer os.Remove(tempFile.Name())
		_, err = tempFile.WriteString(*req.Msg.ParametersSchema)
		if err != nil {
			tempFile.Close()
			return nil, err
		}
		err = tempFile.Close()
		if err != nil {
			return nil, err
		}
		// validate schema
		compiler := jsonschema.NewCompiler()
		schema, err := compiler.Compile(tempFile.Name())
		if err != nil {
			return nil, err
		}
		if schema == nil {
			return nil, fmt.Errorf("Schema is nil")
		}
		parametersSchema = mssql.NVarCharMax(*req.Msg.ParametersSchema)
	}
	if req.Msg.ImageName == "" {
		return nil, fmt.Errorf("ImageName is required")
	}
	inputSpecificationIdBytes, err := conversion.StringToUniqueIdentifier(req.Msg.InputSpecificationId)
	if err != nil {
		return nil, err
	}
	outputSpecificationIdBytes, err := conversion.StringToUniqueIdentifier(req.Msg.OutputSpecificationId)
	if err != nil {
		return nil, err
	}
	database := model.New(middleware.GetTx(ctx))
	workspaceId, err := database.GetWorkspaceIdBySpecificationId(ctx, model.GetWorkspaceIdBySpecificationIdParams{
		SpecificationId: inputSpecificationIdBytes,
	})
	if err != nil {
		return nil, err
	}
	authorization, err := authorization.NewAuthorization(authorization.WithWorkspaceID(workspaceId), authorization.WithWorkspaceRole(api.WorkspaceRole_DEVELOPER)).IsAuthorized(ctx)
	if err != nil {
		return nil, err
	}
	if !authorization {
		return nil, fmt.Errorf("unauthorized")
	}
	cred, err := azidentity.NewClientSecretCredential(s.tenantId, s.clientId, s.clientSecret, nil)
	if err != nil {
		return nil, err
	}
	client, err := azcontainerregistry.NewClient(fmt.Sprintf("https://%s.azurecr.io", s.registryName), cred, nil)
	if err != nil {
		return nil, err
	}
	image, err := client.GetRepositoryProperties(ctx, getScope(workspaceId.String())+"/"+req.Msg.ImageName, nil)
	if err != nil {
		return nil, err
	}
	if image.Name == nil {
		return nil, fmt.Errorf("image not found")
	}
	_, err = database.CreateModel(ctx, model.CreateModelParams{
		InputSpecificationId:  inputSpecificationIdBytes,
		OutputSpecificationId: outputSpecificationIdBytes,
		Name:                  req.Msg.Name,
		ParametersSchema:      &parametersSchema,
		ImageName:             strings.TrimPrefix(*image.Name, getScope(workspaceId.String())+"/"),
	})
	if err != nil {
		return nil, err
	}
	model, err := database.GetModelByInputSpecificationIdAndOutputSpecificationIdAndName(ctx, model.GetModelByInputSpecificationIdAndOutputSpecificationIdAndNameParams{
		InputSpecificationId:  inputSpecificationIdBytes,
		OutputSpecificationId: outputSpecificationIdBytes,
		Name:                  req.Msg.Name,
	})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&api.CreateModelResponse{
		Id: model.ModelID.String(),
	}), nil
}

func (s *Service) GetModels(ctx context.Context, req *connect.Request[api.GetModelsRequest]) (*connect.Response[api.GetModelsResponse], error) {
	if req.Msg.WorkspaceId == "" {
		return nil, fmt.Errorf("WorkspaceId is required")
	}
	workspaceIdBytes, err := conversion.StringToUniqueIdentifier(req.Msg.WorkspaceId)
	if err != nil {
		return nil, err
	}
	database := model.New(middleware.GetTx(ctx))
	models, err := database.GetModelsByWorkspaceId(ctx, model.GetModelsByWorkspaceIdParams{
		WorkspaceId: workspaceIdBytes,
	})
	if err != nil {
		return nil, err
	}
	var modelNames []*api.ModelName
	for _, model := range models {
		modelNames = append(modelNames, &api.ModelName{
			Id:   model.ModelID.String(),
			Name: model.Name,
		})
	}
	return connect.NewResponse(&api.GetModelsResponse{
		Models: modelNames,
	}), nil
}

func (s *Service) GetModel(ctx context.Context, req *connect.Request[api.GetModelRequest]) (*connect.Response[api.GetModelResponse], error) {
	if req.Msg.Id == "" {
		return nil, fmt.Errorf("ModelId is required")
	}
	modelIdBytes, err := conversion.StringToUniqueIdentifier(req.Msg.Id)
	if err != nil {
		return nil, err
	}
	database := model.New(middleware.GetTx(ctx))
	model, err := database.GetModel(ctx, model.GetModelParams{
		ModelId: modelIdBytes,
	})
	if err != nil {
		return nil, err
	}
	var parametersSchema string
	if model.ParametersSchema != nil {
		parametersSchema = string(*model.ParametersSchema)
	}
	return connect.NewResponse(&api.GetModelResponse{
		Model: &api.Model{
			Id:                    model.ModelID.String(),
			InputSpecificationId:  model.InputSpecificationID.String(),
			OutputSpecificationId: model.OutputSpecificationID.String(),
			Name:                  model.Name,
			ParametersSchema:      &parametersSchema,
			ImageName:             model.ImageName,
		},
	}), nil
}

func (s *Service) CreateModelRun(ctx context.Context, req *connect.Request[api.CreateModelRunRequest]) (*connect.Response[api.CreateModelRunResponse], error) {
	if req.Msg.ModelId == "" {
		return nil, fmt.Errorf("ModelId is required")
	}
	if req.Msg.InputFileGroupId == "" {
		return nil, fmt.Errorf("InputFileGroupId is required")
	}
	if req.Msg.Name == "" {
		return nil, fmt.Errorf("Name is required")
	}
	inputFileGroupIdBytes, err := conversion.StringToUniqueIdentifier(req.Msg.InputFileGroupId)
	if err != nil {
		return nil, err
	}
	modelIdBytes, err := conversion.StringToUniqueIdentifier(req.Msg.ModelId)
	if err != nil {
		return nil, err
	}
	database := model.New(middleware.GetTx(ctx))
	appModel, err := database.GetModel(ctx, model.GetModelParams{
		ModelId: modelIdBytes,
	})
	if err != nil {
		return nil, err
	}
	if appModel.ParametersSchema != nil && *appModel.ParametersSchema != "" {
		if req.Msg.Parameters == nil {
			return nil, fmt.Errorf("Parameters are required")
		}
		// create temporary schema file
		tempSchemaFile, err := os.CreateTemp("", "temp")
		if err != nil {
			return nil, err
		}
		defer os.Remove(tempSchemaFile.Name())
		_, err = tempSchemaFile.WriteString(string(*appModel.ParametersSchema))
		if err != nil {
			tempSchemaFile.Close()
			return nil, err
		}
		err = tempSchemaFile.Close()
		if err != nil {
			return nil, err
		}
		compiler := jsonschema.NewCompiler()
		schema, err := compiler.Compile(tempSchemaFile.Name())
		if err != nil {
			return nil, err
		}
		if schema == nil {
			return nil, fmt.Errorf("Schema is nil")
		}
		parameters, err := jsonschema.UnmarshalJSON(strings.NewReader(*req.Msg.Parameters))
		// validate parameters
		err = schema.Validate(parameters)
		if err != nil {
			return nil, err
		}
	}
	_, err = database.CreateModelRun(ctx, model.CreateModelRunParams{
		ModelId:          modelIdBytes,
		InputFileGroupId: inputFileGroupIdBytes,
		Name:             req.Msg.Name,
	})
	if err != nil {
		return nil, err
	}
	modelRun, err := database.GetModelRunByModelIdAndName(ctx, model.GetModelRunByModelIdAndNameParams{
		ModelId: modelIdBytes,
		Name:    req.Msg.Name,
	})
	if err != nil {
		return nil, err
	}
	workspaceId, err := database.GetWorkspaceIdByModelId(ctx, model.GetWorkspaceIdByModelIdParams{
		ModelId: modelIdBytes,
	})
	if err != nil {
		return nil, err
	}
	// upload parameters as json to blob storage
	if appModel.ParametersSchema != nil && *appModel.ParametersSchema != "" {
		cred, err := azidentity.NewClientSecretCredential(s.tenantId, s.clientId, s.clientSecret, nil)
		if err != nil {
			return nil, err
		}
		svcClient, err := service.NewClient(fmt.Sprintf("https://%s.blob.core.windows.net", s.storageAccountName), cred, nil)
		if err != nil {
			return nil, err
		}
		containerClient := svcClient.NewContainerClient(s.modelRunsContainerName)
		blockBlobClient := containerClient.NewBlockBlobClient(getModelRunParametersDirectory(ctx, workspaceId.String(), modelIdBytes.String(), modelRun.ModelRunId.String()) + "/parameters.json")
		_, err = blockBlobClient.UploadBuffer(ctx, []byte(*req.Msg.Parameters), nil)
		if err != nil {
			return nil, err
		}
	}
	config, err := clientcmd.BuildConfigFromFlags("", s.kubeconfigPath)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: strings.ToLower(modelRun.ModelRunId.String()),
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:  appModel.ImageName,
							Image: fmt.Sprintf("%s.azurecr.io/%s/%s:latest", s.registryName, getScope(workspaceId.String()), appModel.ImageName),
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "file-groups-volume",
									MountPath: "/mnt/input",
									SubPath:   getFileGroupDirectory(ctx, workspaceId.String(), appModel.InputSpecificationID.String(), req.Msg.InputFileGroupId),
									ReadOnly:  true,
								},
								{
									Name:      "model-runs-volume",
									MountPath: "/mnt/parameters",
									SubPath:   getModelRunParametersDirectory(ctx, workspaceId.String(), modelIdBytes.String(), modelRun.ModelRunId.String()),
									ReadOnly:  true,
								},
								{
									Name:      "model-runs-volume",
									MountPath: "/mnt/output",
									SubPath:   getModelRunOutputDirectory(ctx, workspaceId.String(), modelIdBytes.String(), modelRun.ModelRunId.String()),
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "file-groups-volume",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "pvc-file-groups-iunvi-dev-eastus-001",
								},
							},
						},
						{
							Name: "model-runs-volume",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "pvc-model-runs-iunvi-dev-eastus-001",
								},
							},
						},
					},
				},
			},
		},
	}
	_, err = clientset.BatchV1().Jobs("default").Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&api.CreateModelRunResponse{
		Id: modelRun.ModelRunId.String(),
	}), nil
}

func (s *Service) GetModelRuns(ctx context.Context, req *connect.Request[api.GetModelRunsRequest]) (*connect.Response[api.GetModelRunsResponse], error) {
	if req.Msg.WorkspaceId == "" {
		return nil, fmt.Errorf("WorkspaceId is required")
	}
	workspaceIdBytes, err := conversion.StringToUniqueIdentifier(req.Msg.WorkspaceId)
	if err != nil {
		return nil, err
	}
	database := model.New(middleware.GetTx(ctx))
	modelRuns, err := database.GetModelRunsByWorkspaceId(ctx, model.GetModelRunsByWorkspaceIdParams{
		WorkspaceId: workspaceIdBytes,
	})
	if err != nil {
		return nil, err
	}
	var apiModelRuns []*api.ModelRun
	for _, modelRun := range modelRuns {
		apiModelRuns = append(apiModelRuns, &api.ModelRun{
			Id:               modelRun.ModelRunId.String(),
			ModelId:          modelRun.ModelID.String(),
			InputFileGroupId: modelRun.InputFileGroupID.String(),
			Name:             modelRun.Name,
		})
	}
	return connect.NewResponse(&api.GetModelRunsResponse{
		ModelRuns: apiModelRuns,
	}), nil
}

func (s *Service) CreateDashboard(ctx context.Context, req *connect.Request[api.CreateDashboardRequest]) (*connect.Response[api.CreateDashboardResponse], error) {
	if req.Msg.ModelId == "" {
		return nil, fmt.Errorf("ModelId is required")
	}
	if req.Msg.Name == "" {
		return nil, fmt.Errorf("Name is required")
	}
	modelIdBytes, err := conversion.StringToUniqueIdentifier(req.Msg.ModelId)
	if err != nil {
		return nil, err
	}
	database := model.New(middleware.GetTx(ctx))
	workspaceId, err := database.GetWorkspaceIdByModelId(ctx, model.GetWorkspaceIdByModelIdParams{
		ModelId: modelIdBytes,
	})
	if err != nil {
		return nil, err
	}
	authorization, err := authorization.NewAuthorization(authorization.WithWorkspaceID(workspaceId), authorization.WithWorkspaceRole(api.WorkspaceRole_DEVELOPER)).IsAuthorized(ctx)
	if err != nil {
		return nil, err
	}
	if !authorization {
		return nil, fmt.Errorf("unauthorized")
	}
	_, err = database.CreateDashboard(ctx, model.CreateDashboardParams{
		ModelId: modelIdBytes,
		Name:    req.Msg.Name,
	})
	if err != nil {
		return nil, err
	}
	dashboard, err := database.GetDashboardByModelIdAndName(ctx, model.GetDashboardByModelIdAndNameParams{
		ModelId: modelIdBytes,
		Name:    req.Msg.Name,
	})
	if err != nil {
		return nil, err
	}
	cred, err := azidentity.NewClientSecretCredential(s.tenantId, s.clientId, s.clientSecret, nil)
	if err != nil {
		return nil, err
	}
	svcClient, err := service.NewClient(fmt.Sprintf("https://%s.blob.core.windows.net", s.storageAccountName), cred, nil)
	if err != nil {
		return nil, err
	}
	containerClient := svcClient.NewContainerClient(s.dashboardsContainerName)
	blockBlobClient := containerClient.NewBlockBlobClient(getDashboardDirectory(ctx, workspaceId.String(), modelIdBytes.String(), dashboard.DashboardID.String()) + "/index.md")
	_, err = blockBlobClient.UploadBuffer(ctx, []byte(req.Msg.Definition), nil)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&api.CreateDashboardResponse{
		Id: dashboard.DashboardID.String(),
	}), nil
}

func (s *Service) GetDashboards(ctx context.Context, req *connect.Request[api.GetDashboardsRequest]) (*connect.Response[api.GetDashboardsResponse], error) {
	if req.Msg.WorkspaceId == "" {
		return nil, fmt.Errorf("WorkspaceId is required")
	}
	workspaceIdBytes, err := conversion.StringToUniqueIdentifier(req.Msg.WorkspaceId)
	if err != nil {
		return nil, err
	}
	var modelIdBytes *mssql.UniqueIdentifier
	if req.Msg.ModelId != nil {
		modelId, err := conversion.StringToUniqueIdentifier(*req.Msg.ModelId)
		if err != nil {
			return nil, err
		}
		modelIdBytes = &modelId
	}
	var modelRunIdBytes *mssql.UniqueIdentifier
	if req.Msg.ModelRunId != nil {
		modelRunId, err := conversion.StringToUniqueIdentifier(*req.Msg.ModelRunId)
		if err != nil {
			return nil, err
		}
		modelRunIdBytes = &modelRunId
	}
	database := model.New(middleware.GetTx(ctx))
	dashboards, err := database.GetDashboardsByWorkspaceIdAndModelId(ctx, model.GetDashboardsByWorkspaceIdAndModelIdParams{
		WorkspaceId: workspaceIdBytes,
		ModelId:     modelIdBytes,
		ModelRunId:  modelRunIdBytes,
	})
	if err != nil {
		return nil, err
	}
	var dashboardNames []*api.Dashboard
	for _, dashboard := range dashboards {
		dashboardNames = append(dashboardNames, &api.Dashboard{
			Id:      dashboard.DashboardID.String(),
			ModelId: dashboard.ModelID.String(),
			Name:    dashboard.Name,
		})
	}
	return connect.NewResponse(&api.GetDashboardsResponse{
		Dashboards: dashboardNames,
	}), nil
}

func (s *Service) GetDashboard(ctx context.Context, req *connect.Request[api.GetDashboardRequest]) (*connect.Response[api.GetDashboardResponse], error) {
	if req.Msg.Id == "" {
		return nil, fmt.Errorf("DashboardId is required")
	}
	dashboardIdBytes, err := conversion.StringToUniqueIdentifier(req.Msg.Id)
	if err != nil {
		return nil, err
	}
	database := model.New(middleware.GetTx(ctx))
	dashboard, err := database.GetDashboard(ctx, model.GetDashboardParams{
		DashboardId: dashboardIdBytes,
	})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&api.GetDashboardResponse{
		Dashboard: &api.Dashboard{
			Id:      dashboard.DashboardID.String(),
			ModelId: dashboard.ModelID.String(),
			Name:    dashboard.Name,
		},
	}), nil
}

func (s *Service) GetModelRunDashboard(ctx context.Context, req *connect.Request[api.GetModelRunDashboardRequest]) (*connect.Response[api.GetModelRunDashboardResponse], error) {
	if req.Msg.ModelRunId == "" {
		return nil, fmt.Errorf("ModelRunId is required")
	}
	if req.Msg.DashboardId == "" {
		return nil, fmt.Errorf("DashboardId is required")
	}
	modelRunIdBytes, err := conversion.StringToUniqueIdentifier(req.Msg.ModelRunId)
	if err != nil {
		return nil, err
	}
	dashboardIdBytes, err := conversion.StringToUniqueIdentifier(req.Msg.DashboardId)
	if err != nil {
		return nil, err
	}
	database := model.New(middleware.GetTx(ctx))
	res, err := database.GetWorkspaceIdAndModelIdByModelRunIdAndDashboardId(ctx, model.GetWorkspaceIdAndModelIdByModelRunIdAndDashboardIdParams{
		ModelRunId:  modelRunIdBytes,
		DashboardId: dashboardIdBytes,
	})
	if err != nil {
		return nil, err
	}
	authorization, err := authorization.NewAuthorization(authorization.WithWorkspaceID(res.WorkspaceID), authorization.WithWorkspaceRole(api.WorkspaceRole_VIEWER)).IsAuthorized(ctx)
	if err != nil {
		return nil, err
	}
	if !authorization {
		return nil, fmt.Errorf("unauthorized")
	}
	// check if model dashboard is created in storage
	cred, err := azidentity.NewClientSecretCredential(s.tenantId, s.clientId, s.clientSecret, nil)
	if err != nil {
		return nil, err
	}
	svcClient, err := service.NewClient(fmt.Sprintf("https://%s.blob.core.windows.net", s.storageAccountName), cred, nil)
	if err != nil {
		return nil, err
	}
	modelRunsDashboardsContainerClient := svcClient.NewContainerClient(s.modelRunDashboardsContainerName)
	blockBlobClient := modelRunsDashboardsContainerClient.NewBlockBlobClient(getModelRunDashboardDirectory(ctx, res.WorkspaceID.String(), res.ModelID.String(), req.Msg.ModelRunId, req.Msg.DashboardId) + "/index.html")
	properties, err := blockBlobClient.GetProperties(ctx, nil)
	if err == nil {
		if properties.ContentLength == nil {
			return nil, fmt.Errorf("ContentLength is nil")
		}
		dashboardHtml := make([]byte, *properties.ContentLength)
		_, err = blockBlobClient.DownloadBuffer(ctx, dashboardHtml, &blob.DownloadBufferOptions{
			Range: blob.HTTPRange{
				Offset: 0,
				Count:  *properties.ContentLength,
			},
		})
		if err != nil {
			return nil, err
		}
		return connect.NewResponse(&api.GetModelRunDashboardResponse{
			DashboardHtml: string(dashboardHtml),
		}), nil
	}
	// if it is not then start a job in aks to create the dashboard
	// list output files in model run directory
	modelRunsContainerClient := svcClient.NewContainerClient(s.modelRunsContainerName)
	pager := modelRunsContainerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Prefix: conversion.ValueToPointer(getModelRunOutputDirectory(ctx, res.WorkspaceID.String(), res.ModelID.String(), req.Msg.ModelRunId) + "/"),
	})
	var outputFiles []string
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, blob := range page.Segment.BlobItems {
			if blob.Name == nil {
				continue
			}
			outputFiles = append(outputFiles, "static/data/parquets/"+strings.TrimPrefix(*blob.Name, getModelRunOutputDirectory(ctx, res.WorkspaceID.String(), res.ModelID.String(), req.Msg.ModelRunId)+"/"))
		}
	}
	// list dashboard files in dashboard directory
	manifestJson := map[string]interface{}{
		"renderedFiles": map[string]interface{}{
			"parquets": outputFiles,
		},
	}
	manifestJsonBytes, err := json.Marshal(manifestJson)
	if err != nil {
		return nil, err
	}
	config, err := clientcmd.BuildConfigFromFlags("", s.kubeconfigPath)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: strings.ToLower(req.Msg.ModelRunId)[:30] + "-" + strings.ToLower(req.Msg.DashboardId)[:30],
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:  "apply-dashboard",
							Image: fmt.Sprintf("%s.azurecr.io/apply-dashboard:latest", s.registryName),
							Command: []string{
								"/bin/sh",
								"-c",
							},
							Args: []string{
								"echo '" + string(manifestJsonBytes) + "' > /app/.evidence/template/static/data/manifest.json && " +
									"echo -e '\n\ndeployment:\n  basePath: /" + s.modelRunDashboardsContainerName + "/" + getModelRunDashboardDirectory(ctx, res.WorkspaceID.String(), res.ModelID.String(), req.Msg.ModelRunId, req.Msg.DashboardId) + "' >> /app/evidence.config.yaml && " +
									"npm run build"},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "model-runs-volume",
									MountPath: "/app/.evidence/template/static/data/parquets",
									SubPath:   getModelRunOutputDirectory(ctx, res.WorkspaceID.String(), res.ModelID.String(), req.Msg.ModelRunId),
									// ReadOnly:  true,
								},
								{
									Name:      "dashboards-volume",
									MountPath: "/app/pages",
									SubPath:   getDashboardDirectory(ctx, res.WorkspaceID.String(), res.ModelID.String(), req.Msg.DashboardId),
									ReadOnly:  true,
								},
								{
									Name:      "web-volume",
									MountPath: "/app/build",
									SubPath:   getModelRunDashboardDirectory(ctx, res.WorkspaceID.String(), res.ModelID.String(), req.Msg.ModelRunId, req.Msg.DashboardId),
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "model-runs-volume",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "pvc-model-runs-iunvi-dev-eastus-001",
								},
							},
						}, {
							Name: "dashboards-volume",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "pvc-dashboards-iunvi-dev-eastus-001",
								},
							},
						}, {
							Name: "web-volume",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "pvc-web-iunvi-dev-eastus-001",
								},
							},
						},
					},
				},
			},
		},
	}
	_, err = clientset.BatchV1().Jobs("default").Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&api.GetModelRunDashboardResponse{
		DashboardHtml: "",
	}), nil
}

func getLandingZoneDirectory(ctx context.Context, workspaceId string) string {
	tenantId := ctx.Value("TenantDirectoryId").(string)
	userObjectId := ctx.Value("UserObjectId").(string)
	return tenantId + "/" + strings.ToLower(workspaceId) + "/" + userObjectId
}

func getFileGroupDirectory(ctx context.Context, workspaceId string, specificationId string, fileGroupId string) string {
	tenantId := ctx.Value("TenantDirectoryId").(string)
	return tenantId + "/" + strings.ToLower(workspaceId) + "/" + strings.ToLower(specificationId) + "/" + strings.ToLower(fileGroupId)
}

func getModelRunOutputDirectory(ctx context.Context, workspaceId string, modelId string, modelRunId string) string {
	tenantId := ctx.Value("TenantDirectoryId").(string)
	return tenantId + "/" + strings.ToLower(workspaceId) + "/" + strings.ToLower(modelId) + "/" + strings.ToLower(modelRunId) + "/output"
}

func getModelRunParametersDirectory(ctx context.Context, workspaceId string, modelId string, modelRunId string) string {
	tenantId := ctx.Value("TenantDirectoryId").(string)
	return tenantId + "/" + strings.ToLower(workspaceId) + "/" + strings.ToLower(modelId) + "/" + strings.ToLower(modelRunId) + "/parameters"
}

func getDashboardDirectory(ctx context.Context, workspaceId string, modelId string, dashboardId string) string {
	tenantId := ctx.Value("TenantDirectoryId").(string)
	return tenantId + "/" + strings.ToLower(workspaceId) + "/" + strings.ToLower(modelId) + "/" + strings.ToLower(dashboardId)
}

func getModelRunDashboardDirectory(ctx context.Context, workspaceId string, modelId string, modelRunId string, dashboardId string) string {
	tenantId := ctx.Value("TenantDirectoryId").(string)
	return "model-run-dashboards/" + tenantId + "/" + strings.ToLower(workspaceId) + "/" + strings.ToLower(modelId) + "/" + strings.ToLower(modelRunId) + "/" + strings.ToLower(dashboardId)
}

func getScope(workspaceId string) string {
	return fmt.Sprintf("scope-%s", strings.ToLower(workspaceId))
}

func (s *Service) getTokenName(workspaceId string) string {
	return fmt.Sprintf("%s-%s", s.registryTokenName, strings.ToLower(workspaceId))
}
