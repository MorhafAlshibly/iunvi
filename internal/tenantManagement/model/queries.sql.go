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
WHERE WorkspaceId = @WorkspaceId;
`

type GetSpecificationsParams struct {
	WorkspaceId mssql.UniqueIdentifier `db:"WorkspaceId"`
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
	rows, err := q.db.QueryContext(ctx, GetSpecifications, sql.Named("WorkspaceId", arg.WorkspaceId))
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
