package dashboard

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime"
	"path/filepath"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/MorhafAlshibly/iunvi/gen/api"
	"github.com/MorhafAlshibly/iunvi/pkg/authorization"
	"github.com/MorhafAlshibly/iunvi/pkg/conversion"
	"github.com/MorhafAlshibly/iunvi/pkg/middleware"
	"github.com/MorhafAlshibly/iunvi/pkg/model"
	"github.com/MorhafAlshibly/iunvi/pkg/sculpt"
	mssql "github.com/microsoft/go-mssqldb"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	blobService "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/sas"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/service"
)

type Service struct {
	tenantId                        string
	clientId                        string
	clientSecret                    string
	registryName                    string
	registryTokenPrefix             string
	storageAccountName              string
	modelRunsContainerName          string
	dashboardsContainerName         string
	modelRunDashboardsContainerName string
	kubeConfigPath                  string
	applyDashboardImageName         string
	modelRunsPVCName                string
	dashboardsPVCName               string
	modelRunDashboardsPVCName       string
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

func WithKubeConfigPath(kubeConfigPath string) func(*Service) {
	return func(input *Service) {
		input.kubeConfigPath = kubeConfigPath
	}
}

func WithApplyDashboardImageName(applyDashboardImageName string) func(*Service) {
	return func(input *Service) {
		input.applyDashboardImageName = applyDashboardImageName
	}
}

func WithModelRunsPVCName(modelRunsPVCName string) func(*Service) {
	return func(input *Service) {
		input.modelRunsPVCName = modelRunsPVCName
	}
}

func WithDashboardsPVCName(dashboardsPVCName string) func(*Service) {
	return func(input *Service) {
		input.dashboardsPVCName = dashboardsPVCName
	}
}

func WithModelRunDashboardsPVCName(modelRunDashboardsPVCName string) func(*Service) {
	return func(input *Service) {
		input.modelRunDashboardsPVCName = modelRunDashboardsPVCName
	}
}

func NewService(options ...func(*Service)) *Service {
	service := &Service{}
	for _, option := range options {
		option(service)
	}
	return service
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
	containerClient := svcClient.NewFileSystemClient(s.dashboardsContainerName)
	fileClient := containerClient.NewFileClient(sculpt.DashboardDirectory(ctx, workspaceId.String(), modelIdBytes.String(), dashboard.DashboardID.String()) + "/index.md")
	err = fileClient.UploadBuffer(ctx, []byte(req.Msg.Definition), nil)
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
	var apiDashboards []*api.Dashboard
	for _, dashboard := range dashboards {
		apiDashboards = append(apiDashboards, &api.Dashboard{
			Id:      dashboard.DashboardID.String(),
			ModelId: dashboard.ModelID.String(),
			Name:    dashboard.Name,
		})
	}
	return connect.NewResponse(&api.GetDashboardsResponse{
		Dashboards: apiDashboards,
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

func (s *Service) GetDashboardMarkdown(ctx context.Context, req *connect.Request[api.GetDashboardRequest]) (*connect.Response[api.GetDashboardMarkdownResponse], error) {
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
	workspaceId, err := database.GetWorkspaceIdByModelId(ctx, model.GetWorkspaceIdByModelIdParams{
		ModelId: dashboard.ModelID,
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
	containerClient := svcClient.NewFileSystemClient(s.dashboardsContainerName)
	fileClient := containerClient.NewFileClient(sculpt.DashboardDirectory(ctx, workspaceId.String(), dashboard.ModelID.String(), dashboard.DashboardID.String()) + "/index.md")
	content, err := fileClient.DownloadStream(ctx, nil)
	if err != nil {
		return nil, err
	}
	markdown := bytes.Buffer{}
	_, err = markdown.ReadFrom(content.Body)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&api.GetDashboardMarkdownResponse{
		Markdown: markdown.String(),
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
	blobSvcClient, err := blobService.NewClient(fmt.Sprintf("https://%s.blob.core.windows.net", s.storageAccountName), cred, nil)
	if err != nil {
		return nil, err
	}
	directory := sculpt.ModelRunDashboardDirectory(ctx, res.WorkspaceID.String(), res.ModelID.String(), req.Msg.ModelRunId, req.Msg.DashboardId)
	// directoryClient := modelRunsDashboardsContainerClient.NewDirectoryClient(directory)
	modelRunDashboardsContainerClient := blobSvcClient.NewContainerClient(s.modelRunDashboardsContainerName)
	indexFileClient := modelRunDashboardsContainerClient.NewBlockBlobClient(directory + "/index.html")
	if err != nil {
		return nil, err
	}
	properties, err := indexFileClient.GetProperties(ctx, nil)
	if err == nil {
		if properties.ContentType == nil {
			// Loop through the files in the directory and set all their mimetypes correctly
			pager := modelRunDashboardsContainerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
				Prefix: conversion.ValueToPointer(directory + "/"),
			})
			for pager.More() {
				page, err := pager.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for _, file := range page.Segment.BlobItems {
					if file.Name == nil {
						continue
					}
					blobClient := modelRunDashboardsContainerClient.NewBlobClient(*file.Name)
					mimeType := mime.TypeByExtension(filepath.Ext(*file.Name))
					_, err = blobClient.SetHTTPHeaders(ctx, blob.HTTPHeaders{
						BlobContentType: conversion.ValueToPointer(mimeType),
					}, nil)
					if err != nil {
						return nil, err
					}
				}
			}
		}
		info := service.KeyInfo{
			Start:  conversion.ValueToPointer(time.Now().UTC().Add(-10 * time.Second).Format(sas.TimeFormat)),
			Expiry: conversion.ValueToPointer(time.Now().UTC().Add(48 * time.Hour).Format(sas.TimeFormat)),
		}
		svcClient, err := service.NewClient(fmt.Sprintf("https://%s.blob.core.windows.net", s.storageAccountName), cred, nil)
		if err != nil {
			return nil, err
		}
		udc, err := svcClient.GetUserDelegationCredential(ctx, info, nil)
		if err != nil {
			return nil, err
		}
		sasQueryParams, err := sas.DatalakeSignatureValues{
			Protocol:       sas.ProtocolHTTPS,
			StartTime:      time.Now().UTC().Add(-10 * time.Second),
			ExpiryTime:     time.Now().UTC().Add(1 * time.Hour),
			Permissions:    (&sas.DirectoryPermissions{Read: true}).String(),
			FileSystemName: s.modelRunDashboardsContainerName,
			DirectoryPath:  directory,
		}.SignWithUserDelegation(udc)
		if err != nil {
			return nil, err
		}
		sasURL := fmt.Sprintf("https://%s.dfs.core.windows.net/%s/%s?%s", s.storageAccountName, s.modelRunDashboardsContainerName, directory, sasQueryParams.Encode())
		return connect.NewResponse(&api.GetModelRunDashboardResponse{
			DashboardSasUrl: sasURL,
		}), nil
	}
	// if it is not then start a job in aks to create the dashboard
	// list output files in model run directory
	modelRunsContainerClient := blobSvcClient.NewContainerClient(s.modelRunsContainerName)
	pager := modelRunsContainerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Prefix: conversion.ValueToPointer(sculpt.ModelRunOutputDirectory(ctx, res.WorkspaceID.String(), res.ModelID.String(), req.Msg.ModelRunId) + "/"),
	})
	var outputFiles []string
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, file := range page.Segment.BlobItems {
			if file.Name == nil {
				continue
			}
			outputFiles = append(outputFiles, "static/data/parquets/"+strings.TrimPrefix(*file.Name, sculpt.ModelRunOutputDirectory(ctx, res.WorkspaceID.String(), res.ModelID.String(), req.Msg.ModelRunId)+"/"))
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
			Name: strings.ToLower(req.Msg.ModelRunId)[:30] + "-" + strings.ToLower(req.Msg.DashboardId)[:30],
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:  s.applyDashboardImageName,
							Image: fmt.Sprintf("%s.azurecr.io/%s:latest", s.registryName, s.applyDashboardImageName),
							Command: []string{
								"/bin/sh",
								"-c",
							},
							Args: []string{
								"echo '" + string(manifestJsonBytes) + "' > /app/.evidence/template/static/data/manifest.json && " +
									"echo -e '\n\ndeployment:\n  basePath: /" + s.modelRunDashboardsContainerName + "/" + sculpt.ModelRunDashboardDirectory(ctx, res.WorkspaceID.String(), res.ModelID.String(), req.Msg.ModelRunId, req.Msg.DashboardId) + "' >> /app/evidence.config.yaml && " +
									// "npm run build",
									"npm run build && " +
									//"mv /app/build/* /app/postbuild && " +
									"echo -e 'self.addEventListener(\"fetch\",(t=>{let e=t.request.url;const s=self.registration.scope;e.startsWith(s+\"entry.html\")&&(e=s+\"index.html\");const n=new URL(e),r=new URLSearchParams(location.href.split(\"?\")[1]);for(const[t,e]of r)n.searchParams.append(t,e);const o=n.toString();t.respondWith(fetch(o))}));' > /app/build/sw.js && " +
									"echo -e '<p>Installing Service Worker, please wait...</p><script>function handleError(n){console.log(n)}navigator.serviceWorker.register(\"sw.js?\"+location.href.split(\"?\")[1]).then((n=>{if(n.installing){const e=n.installing||n.waiting;e.onstatechange=function(){\"installed\"===e.state&&window.location.reload()}}else n.active&&handleError(new Error(\"Service Worker is installed and not redirecting.\"))})).catch(handleError)</script>' > /app/build/entry.html",
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("500m"),
									corev1.ResourceMemory: resource.MustParse("1Gi"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("2000m"),
									corev1.ResourceMemory: resource.MustParse("4Gi"),
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "model-runs-volume",
									MountPath: "/app/.evidence/template/static/data/parquets",
									SubPath:   sculpt.ModelRunOutputDirectory(ctx, res.WorkspaceID.String(), res.ModelID.String(), req.Msg.ModelRunId),
								},
								{
									Name:      "dashboards-volume",
									MountPath: "/app/pages",
									SubPath:   sculpt.DashboardDirectory(ctx, res.WorkspaceID.String(), res.ModelID.String(), req.Msg.DashboardId),
									ReadOnly:  true,
								},
								{
									Name:      "model-run-dashboards-volume",
									MountPath: "/app/build",
									SubPath:   sculpt.ModelRunDashboardDirectory(ctx, res.WorkspaceID.String(), res.ModelID.String(), req.Msg.ModelRunId, req.Msg.DashboardId),
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "model-runs-volume",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: s.modelRunsPVCName,
								},
							},
						}, {
							Name: "dashboards-volume",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: s.dashboardsPVCName,
								},
							},
						}, {
							Name: "model-run-dashboards-volume",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: s.modelRunDashboardsPVCName,
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
		DashboardSasUrl: "",
	}), nil
}
