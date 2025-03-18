package model

import (
	"time"

	mssql "github.com/microsoft/go-mssqldb"
)

type AppFile struct {
	FileID       mssql.UniqueIdentifier `db:"FileId"`
	FileGroupID  mssql.UniqueIdentifier `db:"FileGroupId"`
	FileSchemaID mssql.UniqueIdentifier `db:"FileSchemaId"`
	Name         string                 `db:"Name"`
}

type AppFileGroup struct {
	FileGroupID        mssql.UniqueIdentifier `db:"FileGroupId"`
	SpecificationID    mssql.UniqueIdentifier `db:"SpecificationId"`
	CreatedBy          mssql.UniqueIdentifier `db:"CreatedBy"`
	Name               string                 `db:"Name"`
	ShareWithWorkspace bool                   `db:"ShareWithWorkspace"`
	CreatedAt          time.Time              `db:"CreatedAt"`
}

type AppFileSchema struct {
	FileSchemaID    mssql.UniqueIdentifier `db:"FileSchemaId"`
	SpecificationID mssql.UniqueIdentifier `db:"SpecificationId"`
	FileTypeID      int32                  `db:"FileTypeId"`
	Name            string                 `db:"Name"`
	Definition      mssql.NVarCharMax      `db:"Definition"`
}

type AppModel struct {
	ModelID               mssql.UniqueIdentifier `db:"ModelId"`
	InputSpecificationID  mssql.UniqueIdentifier `db:"InputSpecificationId"`
	OutputSpecificationID mssql.UniqueIdentifier `db:"OutputSpecificationId"`
	Name                  string                 `db:"Name"`
	ParametersSchema      *mssql.NVarCharMax     `db:"ParametersSchema"`
	ImageName             string                 `db:"ImageName"`
	CreatedAt             time.Time              `db:"CreatedAt"`
}

type AppModelRun struct {
	ModelRunID        mssql.UniqueIdentifier `db:"ModelRunId"`
	ModelID           mssql.UniqueIdentifier `db:"ModelId"`
	InputFileGroupID  mssql.UniqueIdentifier `db:"InputFileGroupId"`
	OutputFileGroupID mssql.UniqueIdentifier `db:"OutputFileGroupId"`
	Name              string                 `db:"Name"`
	CreatedAt         mssql.UniqueIdentifier `db:"CreatedAt"`
}

type AppSpecification struct {
	SpecificationID mssql.UniqueIdentifier `db:"SpecificationId"`
	WorkspaceID     mssql.UniqueIdentifier `db:"WorkspaceId"`
	DataModeID      int32                  `db:"DataModeId"`
	Name            string                 `db:"Name"`
	CreatedAt       time.Time              `db:"CreatedAt"`
}

type AuthUserWorkspaceAssignment struct {
	UserObjectID mssql.UniqueIdentifier `db:"UserObjectId"`
	WorkspaceID  mssql.UniqueIdentifier `db:"WorkspaceId"`
	RoleID       int32                  `db:"RoleId"`
}

type AuthWorkspace struct {
	WorkspaceID       mssql.UniqueIdentifier `db:"WorkspaceId"`
	TenantDirectoryID mssql.UniqueIdentifier `db:"TenantDirectoryId"`
	Name              string                 `db:"Name"`
	CreatedAt         time.Time              `db:"CreatedAt"`
}

type AuthWorkspaceRole struct {
	RoleID int32  `db:"RoleId"`
	Name   string `db:"Name"`
}
