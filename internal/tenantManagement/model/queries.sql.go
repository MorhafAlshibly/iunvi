package model

import (
	"context"
	"database/sql"

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
