package model

import (
	"context"
	"database/sql"
	"time"

	mssql "github.com/microsoft/go-mssqldb"
)

const CreateWorkspace = `-- name: CreateWorkspace :execresult
INSERT INTO auth.Workspaces (TenantDirectoryId, Name)
VALUES (auth.fn_GetSessionTenantId(), @Name)
`

func (q *Queries) CreateWorkspace(ctx context.Context, name string) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateWorkspace, sql.Named("Name", name))
}

const GetWorkspaces = `-- name: GetWorkspaces :many
SELECT WorkspaceId,
    TenantDirectoryId,
    Name,
    CreatedAt
FROM auth.Workspaces;
`

type GetWorkspacesParams struct {
	Limit  int32 `db:"limit"`
	Offset int32 `db:"offset"`
}

func (q *Queries) GetWorkspaces(ctx context.Context, arg GetWorkspacesParams) ([]AuthWorkspace, error) {
	rows, err := q.db.QueryContext(ctx, GetWorkspaces, sql.Named("Limit", arg.Limit), sql.Named("Offset", arg.Offset))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AuthWorkspace
	for rows.Next() {
		var i AuthWorkspace
		if err := rows.Scan(
			&i.WorkspaceID,
			&i.TenantDirectoryID,
			&i.Name,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const GetWorkspaceByName = `-- name: GetWorkspaceByName :one
SELECT WorkspaceId,
	TenantDirectoryId,
	Name,
	CreatedAt
FROM auth.Workspaces
WHERE Name = @Name;
`

func (q *Queries) GetWorkspace(ctx context.Context, name string) (AuthWorkspace, error) {
	row := q.db.QueryRowContext(ctx, GetWorkspaceByName, sql.Named("Name", name))
	var item AuthWorkspace
	err := row.Scan(
		&item.WorkspaceID,
		&item.TenantDirectoryID,
		&item.Name,
		&item.CreatedAt,
	)
	return item, err
}

const EditWorkspace = `-- name: EditWorkspace :execresult
UPDATE auth.Workspaces
SET Name = @Name
WHERE WorkspaceId = @WorkspaceId;
`

type EditWorkspaceParams struct {
	WorkspaceId mssql.UniqueIdentifier `db:"WorkspaceId"`
	Name        string                 `db:"Name"`
}

func (q *Queries) EditWorkspace(ctx context.Context, arg EditWorkspaceParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, EditWorkspace, sql.Named("WorkspaceId", arg.WorkspaceId), sql.Named("Name", arg.Name))
}

const GetUserWorkspaceAssignment = `-- name: GetUserWorkspaceAssignment :one
SELECT uwa.RoleId,
       wr.Name AS RoleName
FROM auth.UserWorkspaceAssignments uwa
JOIN auth.WorkspaceRoles wr ON uwa.RoleId = wr.RoleId
WHERE uwa.UserObjectId = @UserObjectId 
  AND uwa.WorkspaceId = @WorkspaceId;
`

type GetUserWorkspaceAssignmentParams struct {
	UserObjectId mssql.UniqueIdentifier `db:"UserObjectId"`
	WorkspaceId  mssql.UniqueIdentifier `db:"WorkspaceId"`
}

type GetUserWorkspaceAssignmentRow struct {
	RoleId   int32
	RoleName string
}

func (q *Queries) GetUserWorkspaceAssignment(ctx context.Context, arg GetUserWorkspaceAssignmentParams) (GetUserWorkspaceAssignmentRow, error) {
	row := q.db.QueryRowContext(ctx, GetUserWorkspaceAssignment, sql.Named("UserObjectId", arg.UserObjectId), sql.Named("WorkspaceId", arg.WorkspaceId))
	var item GetUserWorkspaceAssignmentRow
	err := row.Scan(
		&item.RoleId,
		&item.RoleName,
	)
	return item, err
}

const DeleteUserWorkspaceAssignment = `-- name: DeleteUserWorkspaceAssignment :execresult
DELETE FROM auth.UserWorkspaceAssignments
WHERE UserObjectId = @UserObjectId
  AND WorkspaceId = @WorkspaceId;
`

type DeleteUserWorkspaceAssignmentParams struct {
	UserObjectId mssql.UniqueIdentifier `db:"UserObjectId"`
	WorkspaceId  mssql.UniqueIdentifier `db:"WorkspaceId"`
}

func (q *Queries) DeleteUserWorkspaceAssignment(ctx context.Context, arg DeleteUserWorkspaceAssignmentParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, DeleteUserWorkspaceAssignment, sql.Named("UserObjectId", arg.UserObjectId), sql.Named("WorkspaceId", arg.WorkspaceId))
}

const AssignUserToWorkspace = `-- name: AssignUserToWorkspace :execresult
INSERT INTO auth.UserWorkspaceAssignments (UserObjectId, WorkspaceId, RoleId)
SELECT @UserObjectId, @WorkspaceId, wr.RoleId
FROM auth.WorkspaceRoles wr
WHERE wr.Name = @RoleName;
`

type AssignUserToWorkspaceParams struct {
	UserObjectId mssql.UniqueIdentifier `db:"UserObjectId"`
	WorkspaceId  mssql.UniqueIdentifier `db:"WorkspaceId"`
	RoleName     string                 `db:"RoleName"`
}

func (q *Queries) AssignUserToWorkspace(ctx context.Context, arg AssignUserToWorkspaceParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, AssignUserToWorkspace, sql.Named("UserObjectId", arg.UserObjectId), sql.Named("WorkspaceId", arg.WorkspaceId), sql.Named("RoleName", arg.RoleName))
}

const AuthorizationCheck = `-- name: AuthorizationCheck :one
SELECT COUNT(*)
FROM auth.UserWorkspaceAssignments uwa
JOIN auth.WorkspaceRoles wr ON uwa.RoleId = wr.RoleId
WHERE uwa.UserObjectId = @UserObjectId
  AND uwa.WorkspaceId = @WorkspaceId
  AND wr.Name = @RoleName;
`

type AuthorizationCheckParams struct {
	UserObjectId mssql.UniqueIdentifier `db:"UserObjectId"`
	WorkspaceId  mssql.UniqueIdentifier `db:"WorkspaceId"`
	RoleName     string                 `db:"RoleName"`
}

func (q *Queries) AuthorizationCheck(ctx context.Context, arg AuthorizationCheckParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, AuthorizationCheck, sql.Named("UserObjectId", arg.UserObjectId), sql.Named("WorkspaceId", arg.WorkspaceId), sql.Named("RoleName", arg.RoleName))
	var count int32
	err := row.Scan(&count)
	return count, err
}

const CreateSpecification = `-- name: CreateSpecification :execresult
INSERT INTO app.Specifications (WorkspaceId, DataModeId, Name)
SELECT @WorkspaceId, dm.DataModeId, @Name
FROM app.DataModes dm
WHERE dm.Name = @DataModeName;
`

type CreateSpecificationParams struct {
	WorkspaceId  mssql.UniqueIdentifier `db:"WorkspaceId"`
	DataModeName string                 `db:"DataModeName"`
	Name         string                 `db:"Name"`
}

func (q *Queries) CreateSpecification(ctx context.Context, arg CreateSpecificationParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateSpecification, sql.Named("WorkspaceId", arg.WorkspaceId), sql.Named("DataModeName", arg.DataModeName), sql.Named("Name", arg.Name))
}

const CreateFileSchema = `-- name: CreateFileSchema :execresult
INSERT INTO app.FileSchemas (SpecificationId, FileTypeId, Name, Definition)
SELECT @SpecificationId, ft.FileTypeId, @Name, @Definition
FROM app.FileTypes ft
WHERE ft.Name = @FileTypeName;
`

type CreateFileSchemaParams struct {
	SpecificationId mssql.UniqueIdentifier `db:"SpecificationId"`
	FileTypeName    string                 `db:"FileTypeName"`
	Name            string                 `db:"Name"`
	Definition      string                 `db:"Definition"`
}

func (q *Queries) CreateFileSchema(ctx context.Context, arg CreateFileSchemaParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateFileSchema, sql.Named("SpecificationId", arg.SpecificationId), sql.Named("FileTypeName", arg.FileTypeName), sql.Named("Name", arg.Name), sql.Named("Definition", arg.Definition))
}

const GetSpecificationByWorkspaceIdAndName = `-- name: GetSpecificationByWorkspaceIdAndName :one
SELECT SpecificationId,
	   WorkspaceId,
	   DataModeId,
	   Name,
	   CreatedAt
FROM app.Specifications
WHERE WorkspaceId = @WorkspaceId
  AND Name = @Name;
`

type GetSpecificationByWorkspaceIdAndNameParams struct {
	WorkspaceId mssql.UniqueIdentifier `db:"WorkspaceId"`
	Name        string                 `db:"Name"`
}

func (q *Queries) GetSpecificationByWorkspaceIdAndName(ctx context.Context, arg GetSpecificationByWorkspaceIdAndNameParams) (AppSpecification, error) {
	row := q.db.QueryRowContext(ctx, GetSpecificationByWorkspaceIdAndName, sql.Named("WorkspaceId", arg.WorkspaceId), sql.Named("Name", arg.Name))
	var item AppSpecification
	err := row.Scan(
		&item.SpecificationID,
		&item.WorkspaceID,
		&item.DataModeID,
		&item.Name,
		&item.CreatedAt,
	)
	return item, err
}

const GetSpecifications = `-- name: GetSpecifications :many
SELECT s.SpecificationId,
	   s.WorkspaceId,
	   s.DataModeId,
	   s.Name,
	   s.CreatedAt,
	   dm.Name AS DataModeName
FROM app.Specifications s
JOIN app.DataModes dm ON s.DataModeId = dm.DataModeId
WHERE WorkspaceId = @WorkspaceId
  AND (@DataModeName IS NULL OR dm.Name = @DataModeName);
`

type GetSpecificationsParams struct {
	WorkspaceId  mssql.UniqueIdentifier `db:"WorkspaceId"`
	DataModeName *string                `db:"DataModeName"`
}

type GetSpecificationRow struct {
	SpecificationID mssql.UniqueIdentifier `db:"SpecificationId"`
	WorkspaceID     mssql.UniqueIdentifier `db:"WorkspaceId"`
	DataModeID      int32                  `db:"DataModeId"`
	Name            string                 `db:"Name"`
	CreatedAt       time.Time              `db:"CreatedAt"`
	DataModeName    string                 `db:"DataModeName"`
}

func (q *Queries) GetSpecifications(ctx context.Context, arg GetSpecificationsParams) ([]GetSpecificationRow, error) {
	rows, err := q.db.QueryContext(ctx, GetSpecifications, sql.Named("WorkspaceId", arg.WorkspaceId), sql.Named("DataModeName", arg.DataModeName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSpecificationRow
	for rows.Next() {
		var i GetSpecificationRow
		if err := rows.Scan(
			&i.SpecificationID,
			&i.WorkspaceID,
			&i.DataModeID,
			&i.Name,
			&i.CreatedAt,
			&i.DataModeName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const GetSpecification = `-- name: GetSpecification :one
SELECT s.SpecificationId,
	   s.WorkspaceId,
	   s.DataModeId,
	   s.Name,
	   s.CreatedAt,
	   dm.Name AS DataModeName
FROM app.Specifications s
JOIN app.DataModes dm ON s.DataModeId = dm.DataModeId
WHERE SpecificationId = @SpecificationId;
`

type GetSpecificationParams struct {
	SpecificationId mssql.UniqueIdentifier `db:"SpecificationId"`
}

func (q *Queries) GetSpecification(ctx context.Context, arg GetSpecificationParams) (GetSpecificationRow, error) {
	row := q.db.QueryRowContext(ctx, GetSpecification, sql.Named("SpecificationId", arg.SpecificationId))
	var item GetSpecificationRow
	err := row.Scan(
		&item.SpecificationID,
		&item.WorkspaceID,
		&item.DataModeID,
		&item.Name,
		&item.CreatedAt,
		&item.DataModeName,
	)
	return item, err
}

const GetFileSchemasBySpecificationIdAndDataTypeName = `-- name: GetFileSchemaBySpecificationIdAndDataTypeName :one
SELECT fs.FileSchemaId,
	   fs.SpecificationId,
	   fs.FileTypeId,
	   fs.Name,
	   fs.Definition,
	   ft.Name AS FileTypeName
FROM app.FileSchemas fs
JOIN app.FileTypes ft ON fs.FileTypeId = ft.FileTypeId
WHERE SpecificationId = @SpecificationId
  AND ft.Name = @FileTypeName;
`

type GetFileSchemasBySpecificationIdAndDataTypeNameParams struct {
	SpecificationId mssql.UniqueIdentifier `db:"SpecificationId"`
	FileTypeName    string                 `db:"DataTypeName"`
}

type GetFileSchemaBySpecificationIdAndDataTypeNameRow struct {
	FileSchemaID    mssql.UniqueIdentifier `db:"FileSchemaId"`
	SpecificationID mssql.UniqueIdentifier `db:"SpecificationId"`
	FileTypeID      int32                  `db:"FileTypeId"`
	Name            string                 `db:"Name"`
	Definition      string                 `db:"Definition"`
	FileTypeName    string                 `db:"FileTypeName"`
}

func (q *Queries) GetFileSchemasBySpecificationIdAndDataTypeName(ctx context.Context, arg GetFileSchemasBySpecificationIdAndDataTypeNameParams) ([]GetFileSchemaBySpecificationIdAndDataTypeNameRow, error) {
	rows, err := q.db.QueryContext(ctx, GetFileSchemasBySpecificationIdAndDataTypeName, sql.Named("SpecificationId", arg.SpecificationId), sql.Named("FileTypeName", arg.FileTypeName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFileSchemaBySpecificationIdAndDataTypeNameRow
	for rows.Next() {
		var i GetFileSchemaBySpecificationIdAndDataTypeNameRow
		if err := rows.Scan(
			&i.FileSchemaID,
			&i.SpecificationID,
			&i.FileTypeID,
			&i.Name,
			&i.Definition,
			&i.FileTypeName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const GetWorkspaceIdBySpecificationId = `-- name: GetWorkspaceIdBySpecificationId :one
SELECT WorkspaceId
FROM app.Specifications
WHERE SpecificationId = @SpecificationId;
`

type GetWorkspaceIdBySpecificationIdParams struct {
	SpecificationId mssql.UniqueIdentifier `db:"SpecificationId"`
}

func (q *Queries) GetWorkspaceIdBySpecificationId(ctx context.Context, arg GetWorkspaceIdBySpecificationIdParams) (mssql.UniqueIdentifier, error) {
	row := q.db.QueryRowContext(ctx, GetWorkspaceIdBySpecificationId, sql.Named("SpecificationId", arg.SpecificationId))
	var item mssql.UniqueIdentifier
	err := row.Scan(&item)
	return item, err
}

const CreateFileGroup = `-- name: CreateFileGroup :execresult
INSERT INTO app.FileGroups (SpecificationId, CreatedBy, Name, ShareWithWorkspace)
VALUES (@SpecificationId, @CreatedBy, @Name, @ShareWithWorkspace);
`

type CreateFileGroupParams struct {
	SpecificationId    mssql.UniqueIdentifier `db:"SpecificationId"`
	CreatedBy          mssql.UniqueIdentifier `db:"CreatedBy"`
	Name               string                 `db:"Name"`
	ShareWithWorkspace bool                   `db:"ShareWithWorkspace"`
}

func (q *Queries) CreateFileGroup(ctx context.Context, arg CreateFileGroupParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateFileGroup, sql.Named("SpecificationId", arg.SpecificationId), sql.Named("CreatedBy", arg.CreatedBy), sql.Named("Name", arg.Name), sql.Named("ShareWithWorkspace", arg.ShareWithWorkspace))
}

const GetFileGroupBySpecificationIdAndName = `-- name: GetFileGroupBySpecificationIdAndName :one
SELECT FileGroupId,
	   SpecificationId,
	   CreatedBy,
	   Name,
	   ShareWithWorkspace,
	   CreatedAt
FROM app.FileGroups
WHERE SpecificationId = @SpecificationId
  AND Name = @Name;
`

type GetFileGroupBySpecificationIdAndNameParams struct {
	SpecificationId mssql.UniqueIdentifier `db:"SpecificationId"`
	Name            string                 `db:"Name"`
}

func (q *Queries) GetFileGroupBySpecificationIdAndName(ctx context.Context, arg GetFileGroupBySpecificationIdAndNameParams) (AppFileGroup, error) {
	row := q.db.QueryRowContext(ctx, GetFileGroupBySpecificationIdAndName, sql.Named("SpecificationId", arg.SpecificationId), sql.Named("Name", arg.Name))
	var item AppFileGroup
	err := row.Scan(
		&item.FileGroupID,
		&item.SpecificationID,
		&item.CreatedBy,
		&item.Name,
		&item.ShareWithWorkspace,
		&item.CreatedAt,
	)
	return item, err
}

const GetFileGroupById = `-- name: GetFileGroupById :one
SELECT FileGroupId,
	   SpecificationId,
	   CreatedBy,
	   Name,
	   ShareWithWorkspace,
	   CreatedAt
FROM app.FileGroups
WHERE FileGroupId = @FileGroupId;
`

type GetFileGroupByIdParams struct {
	FileGroupId mssql.UniqueIdentifier `db:"FileGroupId"`
}

func (q *Queries) GetFileGroupById(ctx context.Context, arg GetFileGroupByIdParams) (AppFileGroup, error) {
	row := q.db.QueryRowContext(ctx, GetFileGroupById, sql.Named("FileGroupId", arg.FileGroupId))
	var item AppFileGroup
	err := row.Scan(
		&item.FileGroupID,
		&item.SpecificationID,
		&item.CreatedBy,
		&item.Name,
		&item.ShareWithWorkspace,
		&item.CreatedAt,
	)
	return item, err
}

const CreateFile = `-- name: CreateFile :execresult
INSERT INTO app.Files (FileGroupId, FileSchemaId, Name)
SELECT @FileGroupId, fs.FileSchemaId, @Name
FROM app.FileSchemas fs
WHERE fs.Name = @FileSchemaName;
`

type CreateFileParams struct {
	FileGroupId    mssql.UniqueIdentifier `db:"FileGroupId"`
	FileSchemaName string                 `db:"FileSchemaName"`
	Name           string                 `db:"Name"`
}

func (q *Queries) CreateFile(ctx context.Context, arg CreateFileParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateFile, sql.Named("FileGroupId", arg.FileGroupId), sql.Named("FileSchemaName", arg.FileSchemaName), sql.Named("Name", arg.Name))
}

const GetFileByFileGroupIdAndSchemaName = `-- name: GetFileByFileGroupIdAndSchemaName :one
SELECT FileId,
	   FileGroupId,
	   FileSchemaId,
	   Name
FROM app.Files
JOIN app.FileSchemas fs ON Files.FileSchemaId = fs.FileSchemaId
WHERE FileGroupId = @FileGroupId
  AND fs.Name = @FileSchemaName;
`

type GetFileByFileGroupIdAndSchemaNameParams struct {
	FileGroupId    mssql.UniqueIdentifier `db:"FileGroupId"`
	FileSchemaName string                 `db:"FileSchemaName"`
}

func (q *Queries) GetFileByFileGroupIdAndSchemaName(ctx context.Context, arg GetFileByFileGroupIdAndSchemaNameParams) (AppFile, error) {
	row := q.db.QueryRowContext(ctx, GetFileByFileGroupIdAndSchemaName, sql.Named("FileGroupId", arg.FileGroupId), sql.Named("FileSchemaName", arg.FileSchemaName))
	var item AppFile
	err := row.Scan(
		&item.FileID,
		&item.FileGroupID,
		&item.FileSchemaID,
		&item.Name,
	)
	return item, err
}

const GetFileGroups = `-- name: GetFileGroups :many
SELECT fg.FileGroupId,
	   fg.SpecificationId,
	   fg.CreatedBy,
	   fg.Name,
	   fg.ShareWithWorkspace,
	   fg.CreatedAt
FROM app.FileGroups fg
JOIN app.Specifications s ON fg.SpecificationId = s.SpecificationId
WHERE s.WorkspaceId = @WorkspaceId
  AND (@SpecificationId IS NULL OR fg.SpecificationId = @SpecificationId);
`

type GetFileGroupsParams struct {
	WorkspaceId     mssql.UniqueIdentifier  `db:"WorkspaceId"`
	SpecificationId *mssql.UniqueIdentifier `db:"SpecificationId"`
}

func (q *Queries) GetFileGroups(ctx context.Context, arg GetFileGroupsParams) ([]AppFileGroup, error) {
	rows, err := q.db.QueryContext(ctx, GetFileGroups, sql.Named("WorkspaceId", arg.WorkspaceId), sql.Named("SpecificationId", arg.SpecificationId))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AppFileGroup
	for rows.Next() {
		var i AppFileGroup
		if err := rows.Scan(
			&i.FileGroupID,
			&i.SpecificationID,
			&i.CreatedBy,
			&i.Name,
			&i.ShareWithWorkspace,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const GetFilesByFileGroupId = `-- name: GetFilesByFileGroupId :many
SELECT DISTINCT f.FileId,
	   f.FileGroupId,
	   f.FileSchemaId,
	   f.Name,
	   fs.Name AS FileSchemaName
FROM app.Files f
JOIN app.FileSchemas fs ON f.FileSchemaId = fs.FileSchemaId
WHERE f.FileGroupId = @FileGroupId;
`

type GetFilesByFileGroupIdParams struct {
	FileGroupId mssql.UniqueIdentifier `db:"FileGroupId"`
}

type GetFilesByFileGroupIdRow struct {
	FileID         mssql.UniqueIdentifier `db:"FileId"`
	FileGroupID    mssql.UniqueIdentifier `db:"FileGroupId"`
	FileSchemaID   mssql.UniqueIdentifier `db:"FileSchemaId"`
	Name           string                 `db:"Name"`
	FileSchemaName string                 `db:"FileSchemaName"`
}

func (q *Queries) GetFilesByFileGroupId(ctx context.Context, arg GetFilesByFileGroupIdParams) ([]GetFilesByFileGroupIdRow, error) {
	rows, err := q.db.QueryContext(ctx, GetFilesByFileGroupId, sql.Named("FileGroupId", arg.FileGroupId))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFilesByFileGroupIdRow
	for rows.Next() {
		var i GetFilesByFileGroupIdRow
		if err := rows.Scan(
			&i.FileID,
			&i.FileGroupID,
			&i.FileSchemaID,
			&i.Name,
			&i.FileSchemaName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const CreateModel = `-- name: CreateModel :execresult
INSERT INTO app.Models (InputSpecificationId, OutputSpecificationId, Name, ParametersSchema, ImageName)
VALUES (@InputSpecificationId, @OutputSpecificationId, @Name, @ParametersSchema, @ImageName);
`

type CreateModelParams struct {
	InputSpecificationId  mssql.UniqueIdentifier `db:"InputSpecificationId"`
	OutputSpecificationId mssql.UniqueIdentifier `db:"OutputSpecificationId"`
	Name                  string                 `db:"Name"`
	ParametersSchema      *mssql.NVarCharMax     `db:"ParametersSchema"`
	ImageName             string                 `db:"ImageName"`
}

func (q *Queries) CreateModel(ctx context.Context, arg CreateModelParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateModel, sql.Named("InputSpecificationId", arg.InputSpecificationId), sql.Named("OutputSpecificationId", arg.OutputSpecificationId), sql.Named("Name", arg.Name), sql.Named("ParametersSchema", arg.ParametersSchema), sql.Named("ImageName", arg.ImageName))
}

const GetModelByInputSpecificationIdAndOutputSpecificationIdAndName = `-- name: GetModelByInputSpecificationIdAndOutputSpecificationIdAndName :one
SELECT ModelId,
	   InputSpecificationId,
	   OutputSpecificationId,
	   Name,
	   ParametersSchema,
	   ImageName
FROM app.Models
WHERE InputSpecificationId = @InputSpecificationId
  AND OutputSpecificationId = @OutputSpecificationId
  AND Name = @Name;
`

type GetModelByInputSpecificationIdAndOutputSpecificationIdAndNameParams struct {
	InputSpecificationId  mssql.UniqueIdentifier `db:"InputSpecificationId"`
	OutputSpecificationId mssql.UniqueIdentifier `db:"OutputSpecificationId"`
	Name                  string                 `db:"Name"`
}

func (q *Queries) GetModelByInputSpecificationIdAndOutputSpecificationIdAndName(ctx context.Context, arg GetModelByInputSpecificationIdAndOutputSpecificationIdAndNameParams) (AppModel, error) {
	row := q.db.QueryRowContext(ctx, GetModelByInputSpecificationIdAndOutputSpecificationIdAndName, sql.Named("InputSpecificationId", arg.InputSpecificationId), sql.Named("OutputSpecificationId", arg.OutputSpecificationId), sql.Named("Name", arg.Name))
	var item AppModel
	err := row.Scan(
		&item.ModelID,
		&item.InputSpecificationID,
		&item.OutputSpecificationID,
		&item.Name,
		&item.ParametersSchema,
		&item.ImageName,
	)
	return item, err
}

const GetModelsByWorkspaceId = `-- name: GetModelsByWorkspaceId :many
SELECT m.ModelId,
	   m.Name
FROM app.Models m
JOIN app.Specifications s ON m.InputSpecificationId = s.SpecificationId
WHERE s.WorkspaceId = @WorkspaceId;
`

type GetModelsByWorkspaceIdParams struct {
	WorkspaceId mssql.UniqueIdentifier `db:"WorkspaceId"`
}

type GetModelsByWorkspaceIdRow struct {
	ModelID mssql.UniqueIdentifier `db:"ModelId"`
	Name    string                 `db:"Name"`
}

func (q *Queries) GetModelsByWorkspaceId(ctx context.Context, arg GetModelsByWorkspaceIdParams) ([]GetModelsByWorkspaceIdRow, error) {
	rows, err := q.db.QueryContext(ctx, GetModelsByWorkspaceId, sql.Named("WorkspaceId", arg.WorkspaceId))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetModelsByWorkspaceIdRow
	for rows.Next() {
		var i GetModelsByWorkspaceIdRow
		if err := rows.Scan(
			&i.ModelID,
			&i.Name,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const GetWorkspaceIdByModelId = `-- name: GetWorkspaceIdByModelId :one
SELECT s.WorkspaceId
FROM app.Models m
JOIN app.Specifications s ON m.InputSpecificationId = s.SpecificationId
WHERE m.ModelId = @ModelId;
`

type GetWorkspaceIdByModelIdParams struct {
	ModelId mssql.UniqueIdentifier `db:"ModelId"`
}

func (q *Queries) GetWorkspaceIdByModelId(ctx context.Context, arg GetWorkspaceIdByModelIdParams) (mssql.UniqueIdentifier, error) {
	row := q.db.QueryRowContext(ctx, GetWorkspaceIdByModelId, sql.Named("ModelId", arg.ModelId))
	var item mssql.UniqueIdentifier
	err := row.Scan(&item)
	return item, err
}

const GetModel = `-- name: GetModel :one
SELECT ModelId,
	   InputSpecificationId,
	   OutputSpecificationId,
	   Name,
	   ParametersSchema,
	   ImageName,
	   CreatedAt
FROM app.Models
WHERE ModelId = @ModelId;
`

type GetModelParams struct {
	ModelId mssql.UniqueIdentifier `db:"ModelId"`
}

func (q *Queries) GetModel(ctx context.Context, arg GetModelParams) (AppModel, error) {
	row := q.db.QueryRowContext(ctx, GetModel, sql.Named("ModelId", arg.ModelId))
	var item AppModel
	err := row.Scan(
		&item.ModelID,
		&item.InputSpecificationID,
		&item.OutputSpecificationID,
		&item.Name,
		&item.ParametersSchema,
		&item.ImageName,
		&item.CreatedAt,
	)
	return item, err
}

const CreateModelRun = `-- name: CreateModelRun :execresult
INSERT INTO app.ModelRuns (ModelId, InputFileGroupId, Name)
VALUES (@ModelId, @InputFileGroupId, @Name);
`

type CreateModelRunParams struct {
	ModelId          mssql.UniqueIdentifier `db:"ModelId"`
	InputFileGroupId mssql.UniqueIdentifier `db:"InputFileGroup"`
	Name             string                 `db:"Name"`
}

func (q *Queries) CreateModelRun(ctx context.Context, arg CreateModelRunParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateModelRun, sql.Named("ModelId", arg.ModelId), sql.Named("InputFileGroupId", arg.InputFileGroupId), sql.Named("Name", arg.Name))
}

const GetModelRunByModelIdAndName = `-- name: GetModelRunByModelIdAndName :one
SELECT ModelRunId,
	   ModelId,
	   InputFileGroupId,
	   OutputFileGroupId,
	   Name,
	   CreatedAt
FROM app.ModelRuns
WHERE ModelId = @ModelId
  AND Name = @Name;
`

type GetModelRunByModelIdAndNameParams struct {
	ModelId mssql.UniqueIdentifier `db:"ModelId"`
	Name    string                 `db:"Name"`
}

func (q *Queries) GetModelRunByModelIdAndName(ctx context.Context, arg GetModelRunByModelIdAndNameParams) (AppModelRun, error) {
	row := q.db.QueryRowContext(ctx, GetModelRunByModelIdAndName, sql.Named("ModelId", arg.ModelId), sql.Named("Name", arg.Name))
	var item AppModelRun
	err := row.Scan(
		&item.ModelRunId,
		&item.ModelID,
		&item.InputFileGroupID,
		&item.OutputFileGroupID,
		&item.Name,
		&item.CreatedAt,
	)
	return item, err
}

const GetModelRunsByModelId = `-- name: GetModelRunsByModelId :many
SELECT ModelRunId,
	   ModelId,
	   InputFileGroupId,
	   OutputFileGroupId,
	   Name,
	   CreatedAt
FROM app.ModelRuns
WHERE ModelId = @ModelId;
`

type GetModelRunsByModelIdParams struct {
	ModelId mssql.UniqueIdentifier `db:"ModelId"`
}

func (q *Queries) GetModelRunsByModelId(ctx context.Context, arg GetModelRunsByModelIdParams) ([]AppModelRun, error) {
	rows, err := q.db.QueryContext(ctx, GetModelRunsByModelId, sql.Named("ModelId", arg.ModelId))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AppModelRun
	for rows.Next() {
		var i AppModelRun
		if err := rows.Scan(
			&i.ModelRunId,
			&i.ModelID,
			&i.InputFileGroupID,
			&i.OutputFileGroupID,
			&i.Name,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const GetModelRun = `-- name: GetModelRun :one
SELECT ModelRunId,
	   ModelId,
	   InputFileGroupId,
	   OutputFileGroupId,
	   Name,
	   CreatedAt
FROM app.ModelRuns
WHERE RunId = @RunId;
`

type GetModelRunParams struct {
	ModelRunId mssql.UniqueIdentifier `db:"RunId"`
}

func (q *Queries) GetModelRun(ctx context.Context, arg GetModelRunParams) (AppModelRun, error) {
	row := q.db.QueryRowContext(ctx, GetModelRun, sql.Named("RunId", arg.ModelRunId))
	var item AppModelRun
	err := row.Scan(
		&item.ModelRunId,
		&item.ModelID,
		&item.InputFileGroupID,
		&item.OutputFileGroupID,
		&item.Name,
		&item.CreatedAt,
	)
	return item, err
}

const GetModelRunsByWorkspaceId = `-- name: GetModelRunsByWorkspaceId :many
SELECT mr.ModelRunId,
	   mr.ModelId,
	   mr.InputFileGroupId,
	   mr.OutputFileGroupId,
	   mr.Name,
	   mr.CreatedAt
FROM app.ModelRuns mr
JOIN app.Models m ON mr.ModelId = m.ModelId
JOIN app.Specifications s ON m.InputSpecificationId = s.SpecificationId
WHERE s.WorkspaceId = @WorkspaceId;
`

type GetModelRunsByWorkspaceIdParams struct {
	WorkspaceId mssql.UniqueIdentifier `db:"WorkspaceId"`
}

func (q *Queries) GetModelRunsByWorkspaceId(ctx context.Context, arg GetModelRunsByWorkspaceIdParams) ([]AppModelRun, error) {
	rows, err := q.db.QueryContext(ctx, GetModelRunsByWorkspaceId, sql.Named("WorkspaceId", arg.WorkspaceId))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AppModelRun
	for rows.Next() {
		var i AppModelRun
		if err := rows.Scan(
			&i.ModelRunId,
			&i.ModelID,
			&i.InputFileGroupID,
			&i.OutputFileGroupID,
			&i.Name,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

type CreateDashboardParams struct {
	ModelId mssql.UniqueIdentifier `db:"ModelId"`
	Name    string                 `db:"Name"`
}

const CreateDashboard = `-- name: CreateDashboard :execresult
INSERT INTO app.Dashboards (ModelId, Name)
VALUES (@ModelId, @Name);
`

func (q *Queries) CreateDashboard(ctx context.Context, arg CreateDashboardParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateDashboard, sql.Named("ModelId", arg.ModelId), sql.Named("Name", arg.Name))
}

const GetDashboardByModelIdAndName = `-- name: GetDashboardByModelIdAndName :one
SELECT DashboardId,
	   ModelId,
	   Name,
	   CreatedAt
FROM app.Dashboards
WHERE ModelId = @ModelId
  AND Name = @Name;
`

type GetDashboardByModelIdAndNameParams struct {
	ModelId mssql.UniqueIdentifier `db:"ModelId"`
	Name    string                 `db:"Name"`
}

func (q *Queries) GetDashboardByModelIdAndName(ctx context.Context, arg GetDashboardByModelIdAndNameParams) (AppDashboard, error) {
	row := q.db.QueryRowContext(ctx, GetDashboardByModelIdAndName, sql.Named("ModelId", arg.ModelId), sql.Named("Name", arg.Name))
	var item AppDashboard
	err := row.Scan(
		&item.DashboardID,
		&item.ModelID,
		&item.Name,
		&item.CreatedAt,
	)
	return item, err
}

const GetDashboardsByWorkspaceIdAndModelId = `-- name: GetDashboardsByWorkspaceIdAndModelId :many
SELECT DISTINCT d.DashboardId,
	   d.ModelId,
	   d.Name,
	   d.CreatedAt
FROM app.Dashboards d
JOIN app.Models m ON d.ModelId = m.ModelId
JOIN app.Specifications s ON m.InputSpecificationId = s.SpecificationId
JOIN app.ModelRuns mr ON m.ModelId = mr.ModelId
WHERE s.WorkspaceId = @WorkspaceId
 AND (@ModelId IS NULL OR d.ModelId = @ModelId)
 AND (@ModelRunId IS NULL OR mr.ModelRunId = @ModelRunId);
`

type GetDashboardsByWorkspaceIdAndModelIdParams struct {
	WorkspaceId mssql.UniqueIdentifier  `db:"WorkspaceId"`
	ModelId     *mssql.UniqueIdentifier `db:"ModelId"`
	ModelRunId  *mssql.UniqueIdentifier `db:"ModelRunId"`
}

func (q *Queries) GetDashboardsByWorkspaceIdAndModelId(ctx context.Context, arg GetDashboardsByWorkspaceIdAndModelIdParams) ([]AppDashboard, error) {
	rows, err := q.db.QueryContext(ctx, GetDashboardsByWorkspaceIdAndModelId, sql.Named("WorkspaceId", arg.WorkspaceId), sql.Named("ModelId", arg.ModelId), sql.Named("ModelRunId", arg.ModelRunId))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AppDashboard
	for rows.Next() {
		var i AppDashboard
		if err := rows.Scan(
			&i.DashboardID,
			&i.ModelID,
			&i.Name,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const GetDashboard = `-- name: GetDashboard :one
SELECT DashboardId,
	   ModelId,
	   Name,
	   CreatedAt
FROM app.Dashboards
WHERE DashboardId = @DashboardId;
`

type GetDashboardParams struct {
	DashboardId mssql.UniqueIdentifier `db:"DashboardId"`
}

func (q *Queries) GetDashboard(ctx context.Context, arg GetDashboardParams) (AppDashboard, error) {
	row := q.db.QueryRowContext(ctx, GetDashboard, sql.Named("DashboardId", arg.DashboardId))
	var item AppDashboard
	err := row.Scan(
		&item.DashboardID,
		&item.ModelID,
		&item.Name,
		&item.CreatedAt,
	)
	return item, err
}

const GetWorkspaceIdAndModelIdByModelRunIdAndDashboardId = `-- name: GetWorkspaceIdByModelRunIdAndDashboardId :one
SELECT s.WorkspaceId, mr.ModelId
FROM app.ModelRuns mr
JOIN app.Models m ON mr.ModelId = m.ModelId
JOIN app.Specifications s ON m.InputSpecificationId = s.SpecificationId
JOIN app.Dashboards d ON m.ModelId = d.ModelId
WHERE mr.ModelRunId = @ModelRunId
  AND d.DashboardId = @DashboardId;
`

type GetWorkspaceIdAndModelIdByModelRunIdAndDashboardIdParams struct {
	ModelRunId  mssql.UniqueIdentifier `db:"ModelRunId"`
	DashboardId mssql.UniqueIdentifier `db:"DashboardId"`
}

type GetWorkspaceIdAndModelIdByModelRunIdAndDashboardIdRow struct {
	WorkspaceID mssql.UniqueIdentifier `db:"WorkspaceId"`
	ModelID     mssql.UniqueIdentifier `db:"ModelId"`
}

func (q *Queries) GetWorkspaceIdAndModelIdByModelRunIdAndDashboardId(ctx context.Context, arg GetWorkspaceIdAndModelIdByModelRunIdAndDashboardIdParams) (GetWorkspaceIdAndModelIdByModelRunIdAndDashboardIdRow, error) {
	row := q.db.QueryRowContext(ctx, GetWorkspaceIdAndModelIdByModelRunIdAndDashboardId, sql.Named("ModelRunId", arg.ModelRunId), sql.Named("DashboardId", arg.DashboardId))
	var item GetWorkspaceIdAndModelIdByModelRunIdAndDashboardIdRow
	err := row.Scan(&item.WorkspaceID, &item.ModelID)
	return item, err
}
