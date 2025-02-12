package model

import (
	"context"
	"database/sql"
)

const CreateWorkspace = `-- name: CreateWorkspace :execresult
INSERT INTO auth.Workspaces (TenantDirectoryId, Name)
VALUES (auth.fn_GetSessionTenantId(), @Name)
`

func (q *Queries) CreateWorkspace(ctx context.Context, name string) (sql.Result, error) {
	return q.db.ExecContext(ctx, CreateWorkspace, name)
}

const GetWorkspaces = `-- name: GetWorkspaces :many
SELECT WorkspaceId,
    TenantDirectoryId,
    Name,
    CreatedAt
FROM auth.Workspaces
LIMIT @Limit OFFSET @Offset
`

type GetWorkspacesParams struct {
	Limit  int32 `db:"limit"`
	Offset int32 `db:"offset"`
}

func (q *Queries) GetWorkspaces(ctx context.Context, arg GetWorkspacesParams) ([]AuthWorkspace, error) {
	rows, err := q.db.QueryContext(ctx, GetWorkspaces, arg.Limit, arg.Offset)
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
WHERE Name = @Name
LIMIT 1
`

func (q *Queries) GetWorkspace(ctx context.Context, name string) (AuthWorkspace, error) {
	row := q.db.QueryRowContext(ctx, GetWorkspaceByName, name)
	var item AuthWorkspace
	err := row.Scan(
		&item.WorkspaceID,
		&item.TenantDirectoryID,
		&item.Name,
		&item.CreatedAt,
	)
	return item, err
}
