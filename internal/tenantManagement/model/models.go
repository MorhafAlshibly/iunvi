package model

import (
	"time"

	mssql "github.com/microsoft/go-mssqldb"
)

type AppFile struct {
	FileID       mssql.UniqueIdentifier `db:"file_id"`
	FileGroupID  mssql.UniqueIdentifier `db:"file_group_id"`
	FileSchemaID mssql.UniqueIdentifier `db:"file_schema_id"`
	Name         string                 `db:"name"`
}

type AppFilegroup struct {
	FileGroupID        mssql.UniqueIdentifier `db:"file_group_id"`
	SpecificationID    mssql.UniqueIdentifier `db:"specification_id"`
	CreatedBy          mssql.UniqueIdentifier `db:"created_by"`
	Name               string                 `db:"name"`
	ShareWithWorkspace bool                   `db:"share_with_workspace"`
	CreatedAt          time.Time              `db:"created_at"`
}

type AppFileschema struct {
	FileSchemaID    mssql.UniqueIdentifier `db:"file_schema_id"`
	SpecificationID mssql.UniqueIdentifier `db:"specification_id"`
	Name            string                 `db:"name"`
	Definition      mssql.NVarCharMax      `db:"definition"`
}

type AppModel struct {
	ModelID               mssql.UniqueIdentifier `db:"model_id"`
	InputSpecificationID  mssql.UniqueIdentifier `db:"input_specification_id"`
	OutputSpecificationID mssql.UniqueIdentifier `db:"output_specification_id"`
	Name                  string                 `db:"name"`
	ImageID               string                 `db:"image_id"`
}

type AppModelrun struct {
	RunID             mssql.UniqueIdentifier `db:"run_id"`
	ModelID           mssql.UniqueIdentifier `db:"model_id"`
	StatusID          int32                  `db:"status_id"`
	InputFileGroupID  mssql.UniqueIdentifier `db:"input_file_group_id"`
	OutputFileGroupID mssql.UniqueIdentifier `db:"output_file_group_id"`
	ContainerID       string                 `db:"container_id"`
	CreatedAt         mssql.UniqueIdentifier `db:"created_at"`
}

type AppModelrunstatus struct {
	ModelRunStatusID int32  `db:"model_run_status_id"`
	Name             string `db:"name"`
}

type AppSpecification struct {
	SpecificationID mssql.UniqueIdentifier `db:"specification_id"`
	WorkspaceID     mssql.UniqueIdentifier `db:"workspace_id"`
	Name            string                 `db:"name"`
	CreatedAt       time.Time              `db:"created_at"`
}

type AuthUserworkspaceassignment struct {
	UserObjectID mssql.UniqueIdentifier `db:"user_object_id"`
	WorkspaceID  mssql.UniqueIdentifier `db:"workspace_id"`
	RoleID       int32                  `db:"role_id"`
}

type AuthWorkspace struct {
	WorkspaceID       mssql.UniqueIdentifier `db:"workspace_id"`
	TenantDirectoryID mssql.UniqueIdentifier `db:"tenant_directory_id"`
	Name              string                 `db:"name"`
	CreatedAt         time.Time              `db:"created_at"`
}

type AuthWorkspacerole struct {
	RoleID int32  `db:"role_id"`
	Name   string `db:"name"`
}
