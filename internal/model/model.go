package model

import (
	"context"
	"fmt"
	"os"
	"strings"

	"connectrpc.com/connect"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/containers/azcontainerregistry"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry"
	"github.com/MorhafAlshibly/iunvi/gen/api"
	"github.com/MorhafAlshibly/iunvi/pkg/authorization"
	"github.com/MorhafAlshibly/iunvi/pkg/conversion"
	"github.com/MorhafAlshibly/iunvi/pkg/middleware"
	"github.com/MorhafAlshibly/iunvi/pkg/model"
	"github.com/MorhafAlshibly/iunvi/pkg/sculpt"
	mssql "github.com/microsoft/go-mssqldb"
	"google.golang.org/protobuf/types/known/timestamppb"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/service"
	_ "github.com/marcboeker/go-duckdb"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

type Service struct {
	subscriptionId         string
	resourceGroupName      string
	tenantId               string
	clientId               string
	clientSecret           string
	registryName           string
	registryTokenPrefix    string
	storageAccountName     string
	modelRunsContainerName string
	kubeConfigPath         string
	inputFileMountPath     string
	parametersMountPath    string
	outputFileMountPath    string
	fileGroupsPVCName      string
	modelRunsPVCName       string
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

func WithRegistryTokenPrefix(registryTokenPrefix string) func(*Service) {
	return func(input *Service) {
		input.registryTokenPrefix = registryTokenPrefix
	}
}

func WithStorageAccountName(storageAccountName string) func(*Service) {
	return func(input *Service) {
		input.storageAccountName = storageAccountName
	}
}

func WithModelRunsContainerName(modelRunsContainerName string) func(*Service) {
	return func(input *Service) {
		input.modelRunsContainerName = modelRunsContainerName
	}
}

func WithKubeConfigPath(kubeConfigPath string) func(*Service) {
	return func(input *Service) {
		input.kubeConfigPath = kubeConfigPath
	}
}

func WithInputFileMountPath(inputFileMountPath string) func(*Service) {
	return func(input *Service) {
		input.inputFileMountPath = inputFileMountPath
	}
}

func WithParametersMountPath(parametersMountPath string) func(*Service) {
	return func(input *Service) {
		input.parametersMountPath = parametersMountPath
	}
}

func WithOutputFileMountPath(outputFileMountPath string) func(*Service) {
	return func(input *Service) {
		input.outputFileMountPath = outputFileMountPath
	}
}

func WithFileGroupsPVCName(fileGroupsPVCName string) func(*Service) {
	return func(input *Service) {
		input.fileGroupsPVCName = fileGroupsPVCName
	}
}

func WithModelRunsPVCName(modelRunsPVCName string) func(*Service) {
	return func(input *Service) {
		input.modelRunsPVCName = modelRunsPVCName
	}
}

func NewService(options ...func(*Service)) *Service {
	service := &Service{}
	for _, option := range options {
		option(service)
	}
	return service
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
	tokenName := sculpt.RegistryTokenName(req.Msg.WorkspaceId, s.registryTokenPrefix)
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
	tokenName := sculpt.RegistryTokenName(req.Msg.WorkspaceId, s.registryTokenPrefix)
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
	scope := sculpt.RegistryScope(req.Msg.WorkspaceId)
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
	image, err := client.GetRepositoryProperties(ctx, sculpt.RegistryScope(workspaceId.String())+"/"+req.Msg.ImageName, nil)
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
		ImageName:             strings.TrimPrefix(*image.Name, sculpt.RegistryScope(workspaceId.String())+"/"),
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
	// upload parameters as json to file storage
	if appModel.ParametersSchema != nil && *appModel.ParametersSchema != "" {
		cred, err := azidentity.NewClientSecretCredential(s.tenantId, s.clientId, s.clientSecret, nil)
		if err != nil {
			return nil, err
		}
		svcClient, err := service.NewClient(fmt.Sprintf("https://%s.blob.core.windows.net", s.storageAccountName), cred, nil)
		if err != nil {
			return nil, err
		}
		containerClient := svcClient.NewFileSystemClient(s.modelRunsContainerName)
		fileClient := containerClient.NewFileClient(sculpt.ModelRunParametersDirectory(ctx, workspaceId.String(), modelIdBytes.String(), modelRun.ModelRunId.String()) + "/parameters.json")
		err = fileClient.UploadBuffer(ctx, []byte(*req.Msg.Parameters), nil)
		if err != nil {
			return nil, err
		}
	}
	config, err := clientcmd.BuildConfigFromFlags("", s.kubeConfigPath)
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
							Image: fmt.Sprintf("%s.azurecr.io/%s/%s:latest", s.registryName, sculpt.RegistryScope(workspaceId.String()), appModel.ImageName),
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "file-groups-volume",
									MountPath: s.inputFileMountPath,
									SubPath:   sculpt.FileGroupDirectory(ctx, workspaceId.String(), appModel.InputSpecificationID.String(), req.Msg.InputFileGroupId),
									ReadOnly:  true,
								},
								{
									Name:      "model-runs-volume",
									MountPath: s.parametersMountPath,
									SubPath:   sculpt.ModelRunParametersDirectory(ctx, workspaceId.String(), modelIdBytes.String(), modelRun.ModelRunId.String()),
									ReadOnly:  true,
								},
								{
									Name:      "model-runs-volume",
									MountPath: s.outputFileMountPath,
									SubPath:   sculpt.ModelRunOutputDirectory(ctx, workspaceId.String(), modelIdBytes.String(), modelRun.ModelRunId.String()),
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "file-groups-volume",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: s.fileGroupsPVCName,
								},
							},
						},
						{
							Name: "model-runs-volume",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: s.modelRunsPVCName,
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
