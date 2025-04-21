package file

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
	"github.com/MorhafAlshibly/iunvi/gen/api"
	"github.com/MorhafAlshibly/iunvi/pkg/authorization"
	"github.com/MorhafAlshibly/iunvi/pkg/conversion"
	"github.com/MorhafAlshibly/iunvi/pkg/middleware"
	"github.com/MorhafAlshibly/iunvi/pkg/model"
	"github.com/MorhafAlshibly/iunvi/pkg/sculpt"
	mssql "github.com/microsoft/go-mssqldb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	blobService "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/sas"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/service"
	_ "github.com/marcboeker/go-duckdb"
)

type Service struct {
	tenantId                 string
	clientId                 string
	clientSecret             string
	storageAccountName       string
	landingZoneContainerName string
	fileGroupsContainerName  string
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

func NewService(options ...func(*Service)) *Service {
	service := &Service{}
	for _, option := range options {
		option(service)
	}
	return service
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
	directory := sculpt.LandingZoneDirectory(ctx, req.Msg.WorkspaceId)
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
	// Create File Signature Values with desired permissions and sign with user delegation credential
	sasQueryParams, err := sas.DatalakeSignatureValues{
		Protocol:       sas.ProtocolHTTPS,
		StartTime:      time.Now().UTC().Add(-10 * time.Second),
		ExpiryTime:     time.Now().UTC().Add(1 * time.Hour),
		Permissions:    (&sas.FilePermissions{Write: true}).String(),
		FileSystemName: s.landingZoneContainerName,
		FilePath:       directory + "/" + req.Msg.FileName,
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
	directory := sculpt.LandingZoneDirectory(ctx, req.Msg.WorkspaceId)
	svcClient, err := blobService.NewClient(fmt.Sprintf("https://%s.blob.core.windows.net", s.storageAccountName), cred, nil)
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
		for _, file := range page.Segment.BlobItems {
			if file.Name == nil {
				continue
			}
			fileName := strings.TrimPrefix(*file.Name, directory+"/")
			files = append(files, &api.LandingZoneFile{
				Name:         fileName,
				Size:         uint64(conversion.PointerToValue(file.Properties.ContentLength, 0)),
				LastModified: timestamppb.New(conversion.PointerToValue(file.Properties.LastModified, time.Time{})),
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
	svcClient, err := blobService.NewClient(fmt.Sprintf("https://%s.blob.core.windows.net", s.storageAccountName), cred, nil)
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
		blobClient := landingZoneContainerClient.NewBlobClient(sculpt.LandingZoneDirectory(ctx, workspaceId.String()) + "/" + schemaFileMapping.LandingZoneFileName)
		// TODO: Acquire a lease on the file to prevent other users from writing to it

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
		landingZoneBlobClient := landingZoneContainerClient.NewBlobClient(sculpt.LandingZoneDirectory(ctx, workspaceId.String()) + "/" + schemaFileMapping.LandingZoneFileName)
		fileGroupsBlobClient := fileGroupsContainerClient.NewBlobClient(sculpt.FileGroupDirectory(ctx, workspaceId.String(), req.Msg.SpecificationId, fileGroup.FileGroupID.String()) + "/" + schemaFileMapping.LandingZoneFileName)
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
