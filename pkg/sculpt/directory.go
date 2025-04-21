package sculpt

import (
	"context"
	"strings"
)

func LandingZoneDirectory(ctx context.Context, workspaceId string) string {
	tenantId := ctx.Value("TenantDirectoryId").(string)
	userObjectId := ctx.Value("UserObjectId").(string)
	return tenantId + "/" + strings.ToLower(workspaceId) + "/" + userObjectId
}

func FileGroupDirectory(ctx context.Context, workspaceId string, specificationId string, fileGroupId string) string {
	tenantId := ctx.Value("TenantDirectoryId").(string)
	return tenantId + "/" + strings.ToLower(workspaceId) + "/" + strings.ToLower(specificationId) + "/" + strings.ToLower(fileGroupId)
}

func ModelRunOutputDirectory(ctx context.Context, workspaceId string, modelId string, modelRunId string) string {
	tenantId := ctx.Value("TenantDirectoryId").(string)
	return tenantId + "/" + strings.ToLower(workspaceId) + "/" + strings.ToLower(modelId) + "/" + strings.ToLower(modelRunId) + "/output"
}

func ModelRunParametersDirectory(ctx context.Context, workspaceId string, modelId string, modelRunId string) string {
	tenantId := ctx.Value("TenantDirectoryId").(string)
	return tenantId + "/" + strings.ToLower(workspaceId) + "/" + strings.ToLower(modelId) + "/" + strings.ToLower(modelRunId) + "/parameters"
}

func DashboardDirectory(ctx context.Context, workspaceId string, modelId string, dashboardId string) string {
	tenantId := ctx.Value("TenantDirectoryId").(string)
	return tenantId + "/" + strings.ToLower(workspaceId) + "/" + strings.ToLower(modelId) + "/" + strings.ToLower(dashboardId)
}

func ModelRunDashboardDirectory(ctx context.Context, workspaceId string, modelId string, modelRunId string, dashboardId string) string {
	tenantId := ctx.Value("TenantDirectoryId").(string)
	return tenantId + "/" + strings.ToLower(workspaceId) + "/" + strings.ToLower(modelId) + "/" + strings.ToLower(modelRunId) + "/" + strings.ToLower(dashboardId)
}
