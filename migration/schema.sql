-- =============================================
-- Schema: auth (Authentication & Tenancy)
-- =============================================
CREATE SCHEMA auth;
GO;
-- Workspaces (Isolated environments within a tenant)
CREATE TABLE auth.Workspaces (
    WorkspaceId UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    TenantDirectoryId UNIQUEIDENTIFIER NOT NULL,
    Name NVARCHAR(255) NOT NULL,
    CreatedAt DATETIME DEFAULT GETUTCDATE(),
    CONSTRAINT UQ_authWorkspaces_TenantDirectoryId_Name UNIQUE (TenantDirectoryId, Name)
);
-- Workspace Roles (e.g., Developer, User, Viewer)
CREATE TABLE auth.WorkspaceRoles (
    RoleId INT PRIMARY KEY,
    Name NVARCHAR(255) NOT NULL,
    CONSTRAINT UQ_authWorkspaceRoles_Name UNIQUE(Name)
);
INSERT INTO auth.WorkspaceRoles (RoleId, Name)
VALUES (1, 'Developer'),
    (2, 'User'),
    (3, 'Viewer');
-- User access to workspaces
CREATE TABLE auth.UserWorkspaceAssignments (
    UserObjectId UNIQUEIDENTIFIER NOT NULL,
    WorkspaceId UNIQUEIDENTIFIER NOT NULL,
    RoleId INT NOT NULL,
    PRIMARY KEY (UserObjectId, WorkspaceId),
    FOREIGN KEY (WorkspaceId) REFERENCES auth.Workspaces(WorkspaceId),
    FOREIGN KEY (RoleId) REFERENCES auth.WorkspaceRoles(RoleId)
);
-- =============================================
-- Schema: app (Business Functionality)
-- =============================================
GO;
CREATE SCHEMA app;
GO;
-- Data Modes (e.g., Input, Output)
CREATE TABLE app.DataModes (
    DataModeId INT PRIMARY KEY,
    Name NVARCHAR(255) NOT NULL,
    CONSTRAINT UQ_appDataModes_Name UNIQUE(Name)
);
INSERT INTO app.DataModes (DataModeId, Name)
VALUES (1, 'Input'),
    (2, 'Output');
-- File Types (e.g., CSV, JSON, Parquet)
CREATE TABLE app.FileTypes (
    FileTypeId INT PRIMARY KEY,
    DataModeId INT NOT NULL,
    Name NVARCHAR(255) NOT NULL,
    Extension NVARCHAR(10) NOT NULL,
    FOREIGN KEY (DataModeId) REFERENCES app.DataModes(DataModeId),
    CONSTRAINT UQ_appFileTypes_Name UNIQUE(Name)
);
INSERT INTO app.FileTypes (FileTypeId, DataModeId, Name, Extension)
VALUES (1, 1, 'CSV', 'csv'),
    (2, 2, 'Parquet', 'parquet');
-- Specifications (Input/Output schemas)
CREATE TABLE app.Specifications (
    SpecificationId UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    WorkspaceId UNIQUEIDENTIFIER NOT NULL,
    DataModeId INT NOT NULL,
    Name NVARCHAR(255) NOT NULL,
    CreatedAt DATETIME DEFAULT GETUTCDATE(),
    FOREIGN KEY (WorkspaceId) REFERENCES auth.Workspaces(WorkspaceId),
    FOREIGN KEY (DataModeId) REFERENCES app.DataModes(DataModeId),
    CONSTRAINT UQ_authSpecifications_WorkspaceId_Name UNIQUE(WorkspaceId, Name)
);
-- File Schemas (Structure of input files)
CREATE TABLE app.FileSchemas (
    FileSchemaId UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    SpecificationId UNIQUEIDENTIFIER NOT NULL,
    FileTypeId INT NOT NULL,
    Name NVARCHAR(255) NOT NULL,
    Definition NVARCHAR(MAX) NOT NULL,
    FOREIGN KEY (SpecificationId) REFERENCES app.Specifications(SpecificationId),
    FOREIGN KEY (FileTypeId) REFERENCES app.FileTypes(FileTypeId),
    CONSTRAINT UQ_appFileSchemas_SpecificationId_Name UNIQUE(SpecificationId, Name)
);
-- File Groups (Collections of files for a model run)
CREATE TABLE app.FileGroups (
    FileGroupId UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    SpecificationId UNIQUEIDENTIFIER NOT NULL,
    CreatedBy UNIQUEIDENTIFIER NOT NULL,
    Name NVARCHAR(255) NOT NULL,
    ShareWithWorkspace BIT NOT NULL,
    CreatedAt DATETIME DEFAULT GETUTCDATE(),
    FOREIGN KEY (SpecificationId) REFERENCES app.Specifications(SpecificationId),
    CONSTRAINT UQ_appFileGroups_SpecificationId_Name UNIQUE(SpecificationId, Name)
);
-- Files (Individual files within a group)
CREATE TABLE app.Files (
    FileId UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    FileGroupId UNIQUEIDENTIFIER NOT NULL,
    FileSchemaId UNIQUEIDENTIFIER NOT NULL,
    Name NVARCHAR(255) NOT NULL,
    FOREIGN KEY (FileGroupId) REFERENCES app.FileGroups(FileGroupId),
    FOREIGN KEY (FileSchemaId) REFERENCES app.FileSchemas(FileSchemaId),
    CONSTRAINT UQ_appFiles_FileGroupId_FileSchemaId UNIQUE (FileGroupId, FileSchemaId),
    CONSTRAINT UQ_appFiles_FileGroupId_Name UNIQUE (FileGroupId, Name)
);
-- Models (ML/Analytics models)
CREATE TABLE app.Models (
    ModelId UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    InputSpecificationId UNIQUEIDENTIFIER NOT NULL,
    OutputSpecificationId UNIQUEIDENTIFIER NOT NULL,
    Name NVARCHAR(255) NOT NULL,
    ParametersSchema NVARCHAR(MAX) NOT NULL,
    ImageId NVARCHAR(255) NOT NULL,
    FOREIGN KEY (InputSpecificationId) REFERENCES app.Specifications(SpecificationId),
    FOREIGN KEY (OutputSpecificationId) REFERENCES app.Specifications(SpecificationId),
    CONSTRAINT UQ_appModels_InputSpecificationId_OutputSpecificationId_Name UNIQUE(
        InputSpecificationId,
        OutputSpecificationId,
        Name
    )
);
CREATE TABLE app.ModelRunStatuses (
    ModelRunStatusId INT PRIMARY KEY,
    Name NVARCHAR(255) NOT NULL,
    CONSTRAINT UQ_appModelRunStatuses_Name UNIQUE(Name)
);
INSERT INTO app.ModelRunStatuses (ModelRunStatusId, Name)
VALUES (1, 'Pending'),
    (2, 'Running'),
    (3, 'Completed'),
    (4, 'Failed');
-- Model Runs (Executions of a model)
CREATE TABLE app.ModelRuns (
    RunId UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    ModelId UNIQUEIDENTIFIER NOT NULL,
    StatusId INT NOT NULL,
    InputFileGroupId UNIQUEIDENTIFIER NOT NULL,
    OutputFileGroupId UNIQUEIDENTIFIER NOT NULL,
    Parameters NVARCHAR(MAX) NOT NULL,
    ContainerId NVARCHAR(255) NULL,
    CreatedAt DATETIME DEFAULT GETUTCDATE(),
    FOREIGN KEY (ModelId) REFERENCES app.Models(ModelId),
    FOREIGN KEY (StatusId) REFERENCES app.ModelRunStatuses(ModelRunStatusId),
    FOREIGN KEY (InputFileGroupId) REFERENCES app.FileGroups(FileGroupId),
    FOREIGN KEY (OutputFileGroupId) REFERENCES app.FileGroups(FileGroupId)
);
-- =============================================
-- Triggers
-- =============================================
-- Trigger to make sure FileSchema FileTypeId is of the same DataMode as the Specification
GO;
CREATE TRIGGER app.FileSchemas_DataModeCheck ON app.FileSchemas
AFTER
INSERT,
    UPDATE AS BEGIN IF EXISTS (
        SELECT 1
        FROM inserted I
            JOIN app.Specifications S ON I.SpecificationId = S.SpecificationId
            JOIN app.FileTypes FT ON I.FileTypeId = FT.FileTypeId
        WHERE S.DataModeId <> FT.DataModeId
    ) BEGIN RAISERROR(
        'File schema data mode must match the specification data mode',
        16,
        1
    );
ROLLBACK TRANSACTION;
END;
END;
-- Trigger to make sure InputSpecificationId and OutputSpecificationId are from the same workspace
GO;
CREATE TRIGGER app.Models_WorkspaceCheck ON app.Models
AFTER
INSERT,
    UPDATE AS BEGIN IF EXISTS (
        SELECT 1
        FROM inserted I
            JOIN app.Specifications S1 ON I.InputSpecificationId = S1.SpecificationId
            JOIN app.Specifications S2 ON I.OutputSpecificationId = S2.SpecificationId
            JOIN auth.Workspaces W1 ON S1.WorkspaceId = W1.WorkspaceId
            JOIN auth.Workspaces W2 ON S2.WorkspaceId = W2.WorkspaceId
        WHERE W1.WorkspaceId <> W2.WorkspaceId
    ) BEGIN RAISERROR(
        'Input and output specifications must belong to the same workspace',
        16,
        1
    );
ROLLBACK TRANSACTION;
END;
END;
-- Trigger to make sure InputSpecificationId is of DataMode Input and OutputSpecificationId is of DataMode Output
GO;
CREATE TRIGGER app.Models_DataModeCheck ON app.Models
AFTER
INSERT,
    UPDATE AS BEGIN IF EXISTS (
        SELECT 1
        FROM inserted I
            JOIN app.Specifications S1 ON I.InputSpecificationId = S1.SpecificationId
            JOIN app.Specifications S2 ON I.OutputSpecificationId = S2.SpecificationId
        WHERE S1.DataModeId <> 1
            OR S2.DataModeId <> 2
    ) BEGIN RAISERROR(
        'Input specification must be of data mode Input and output specification must be of data mode Output',
        16,
        1
    );
ROLLBACK TRANSACTION;
END;
END;
-- Trigger to make sure InputFileGroupId are of the InputSpecificationId and OutputFileGroupId are of the OutputSpecificationId
GO;
CREATE TRIGGER app.ModelRuns_SpecificationCheck ON app.ModelRuns
AFTER
INSERT,
    UPDATE AS BEGIN IF EXISTS (
        SELECT 1
        FROM inserted I
            JOIN app.Models M ON I.ModelId = M.ModelId
            JOIN app.FileGroups FG1 ON I.InputFileGroupId = FG1.FileGroupId
            JOIN app.FileGroups FG2 ON I.OutputFileGroupId = FG2.FileGroupId
        WHERE FG1.SpecificationId <> M.InputSpecificationId
            OR FG2.SpecificationId <> M.OutputSpecificationId
    ) BEGIN RAISERROR(
        'Input and output file groups must belong to the input and output specifications of the model',
        16,
        1
    );
ROLLBACK TRANSACTION;
END;
END;
GO;
-- =============================================
-- Helper functions
-- =============================================
GO;
CREATE FUNCTION auth.fn_GetSessionTenantId() RETURNS UNIQUEIDENTIFIER WITH SCHEMABINDING AS BEGIN RETURN CAST(
    SESSION_CONTEXT(N'TenantDirectoryId') AS UNIQUEIDENTIFIER
);
END;
GO;
CREATE FUNCTION auth.fn_GetSessionUserId() RETURNS UNIQUEIDENTIFIER WITH SCHEMABINDING AS BEGIN RETURN CAST(
    SESSION_CONTEXT(N'UserObjectId') AS UNIQUEIDENTIFIER
);
END;
GO;
-- =============================================
-- Roles and Permissions
-- =============================================
-- Create roles
CREATE ROLE WebApp;
-- Grant permissions to roles
GRANT SELECT,
    INSERT,
    UPDATE,
    DELETE ON auth.Workspaces TO WebApp;
DENY
UPDATE ON auth.Workspaces(WorkspaceId, TenantDirectoryId, CreatedAt) TO WebApp;
GRANT SELECT ON auth.WorkspaceRoles TO WebApp;
GRANT SELECT,
    INSERT,
    UPDATE,
    DELETE ON auth.UserWorkspaceAssignments TO WebApp;
DENY
UPDATE ON auth.UserWorkspaceAssignments(UserObjectId, WorkspaceId) TO WebApp;
GRANT SELECT ON app.DataModes TO WebApp;
GRANT SELECT ON app.FileTypes TO WebApp;
GRANT SELECT,
    INSERT,
    UPDATE,
    DELETE ON app.Specifications TO WebApp;
DENY
UPDATE ON app.Specifications(
        SpecificationId,
        WorkspaceId,
        DataModeId,
        CreatedAt
    ) TO WebApp;
GRANT SELECT,
    INSERT,
    UPDATE,
    DELETE ON app.FileSchemas TO WebApp;
DENY
UPDATE ON app.FileSchemas(
        FileSchemaId,
        SpecificationId,
        FileTypeId,
        Definition
    ) TO WebApp;
GRANT SELECT,
    INSERT,
    UPDATE,
    DELETE ON app.FileGroups TO WebApp;
DENY
UPDATE ON app.FileGroups(
        FileGroupId,
        SpecificationId,
        CreatedBy,
        CreatedAt
    ) TO WebApp;
GRANT SELECT,
    INSERT,
    UPDATE,
    DELETE ON app.Files TO WebApp;
DENY
UPDATE ON app.Files(FileId, FileGroupId, FileSchemaId) TO WebApp;
GRANT SELECT,
    INSERT,
    UPDATE,
    DELETE ON app.Models TO WebApp;
DENY
UPDATE ON app.Models(
        ModelId,
        InputSpecificationId,
        OutputSpecificationId,
        ParametersSchema,
        ImageId
    ) TO WebApp;
GRANT SELECT ON app.ModelRunStatuses TO WebApp;
GRANT SELECT,
    INSERT,
    DELETE ON app.ModelRuns TO WebApp;
DENY
UPDATE ON app.ModelRuns(
        RunId,
        ModelId,
        StatusId,
        InputFileGroupId,
        OutputFileGroupId,
        Parameters,
        CreatedAt
    ) TO WebApp;
GRANT EXECUTE ON auth.fn_GetSessionTenantId TO WebApp;
GRANT EXECUTE ON auth.fn_GetSessionUserId TO WebApp;
-- =============================================
-- Row-Level Security (RLS) Policies
-- =============================================
-- In the session context I will send the TenantDirectoryId, UserObjectId And TenantRole which is either "Admin" or "User"
-- =============================================
GO;
CREATE FUNCTION auth.fn_Workspaces_Filter(
    @TenantDirectoryId UNIQUEIDENTIFIER,
    @WorkspaceId UNIQUEIDENTIFIER
) RETURNS TABLE WITH SCHEMABINDING AS RETURN
SELECT 1 AS Result
WHERE auth.fn_GetSessionTenantId() = @TenantDirectoryId
    AND (
        SESSION_CONTEXT(N'TenantRole') = N'Admin'
        OR EXISTS(
            SELECT 1
            FROM auth.UserWorkspaceAssignments
            WHERE UserObjectId = auth.fn_GetSessionUserId()
                AND WorkspaceId = @WorkspaceId
        )
    );
GO;
CREATE FUNCTION auth.fn_Workspaces_Block(@TenantDirectoryId UNIQUEIDENTIFIER) RETURNS TABLE WITH SCHEMABINDING AS RETURN
SELECT 1 AS Result
WHERE auth.fn_GetSessionTenantId() = @TenantDirectoryId
    AND SESSION_CONTEXT(N'TenantRole') = N'Admin';
GO;
CREATE SECURITY POLICY auth.WorkspacesPolicy
ADD FILTER PREDICATE auth.fn_Workspaces_Filter(TenantDirectoryId, WorkspaceId) ON auth.Workspaces,
    ADD BLOCK PREDICATE auth.fn_Workspaces_Block(TenantDirectoryId) ON auth.Workspaces WITH (STATE = ON);
-- =============================================
GO;
CREATE FUNCTION auth.UserWorkspaceAssignments_Filter(
    @UserObjectId UNIQUEIDENTIFIER,
    @WorkspaceId UNIQUEIDENTIFIER
) RETURNS TABLE WITH SCHEMABINDING AS RETURN
SELECT 1 AS Result
WHERE auth.fn_GetSessionUserId() = @UserObjectId
    OR (
        SESSION_CONTEXT(N'TenantRole') = N'Admin'
        AND EXISTS(
            SELECT 1
            FROM auth.Workspaces W
            WHERE W.WorkspaceId = @WorkspaceId
                AND W.TenantDirectoryId = auth.fn_GetSessionTenantId()
        )
    );
GO;
CREATE FUNCTION auth.UserWorkspaceAssignments_Block(@WorkspaceId UNIQUEIDENTIFIER) RETURNS TABLE WITH SCHEMABINDING AS RETURN
SELECT 1 AS Result
WHERE SESSION_CONTEXT(N'TenantRole') = N'Admin'
    AND EXISTS(
        SELECT 1
        FROM auth.Workspaces W
        WHERE W.WorkspaceId = @WorkspaceId
            AND W.TenantDirectoryId = auth.fn_GetSessionTenantId()
    );
GO;
CREATE SECURITY POLICY auth.UserWorkspaceAssignmentsPolicy
ADD FILTER PREDICATE auth.UserWorkspaceAssignments_Filter(UserObjectId, WorkspaceId) ON auth.UserWorkspaceAssignments,
    ADD BLOCK PREDICATE auth.UserWorkspaceAssignments_Block(WorkspaceId) ON auth.UserWorkspaceAssignments WITH (STATE = ON);
-- =============================================
GO;
CREATE FUNCTION app.fn_Specifications_Filter(@WorkspaceId UNIQUEIDENTIFIER) RETURNS TABLE WITH SCHEMABINDING AS RETURN
SELECT 1 AS Result
WHERE EXISTS(
        SELECT 1
        FROM auth.UserWorkspaceAssignments
        WHERE UserObjectId = auth.fn_GetSessionUserId()
            AND WorkspaceId = @WorkspaceId
    )
    OR (
        SESSION_CONTEXT(N'TenantRole') = N'Admin'
        AND EXISTS(
            SELECT 1
            FROM auth.Workspaces
            WHERE WorkspaceId = @WorkspaceId
                AND TenantDirectoryId = auth.fn_GetSessionTenantId()
        )
    );
-- If user is admin then spec has to exist in tenant, if not admin then spec has to be assigned to user
GO;
CREATE FUNCTION app.fn_Specifications_Block(@WorkspaceId UNIQUEIDENTIFIER) RETURNS TABLE WITH SCHEMABINDING AS RETURN
SELECT 1 AS Result
WHERE EXISTS(
        SELECT 1
        FROM auth.UserWorkspaceAssignments
        WHERE UserObjectId = auth.fn_GetSessionUserId()
            AND WorkspaceId = @WorkspaceId
            AND RoleId = 1
    )
    OR (
        SESSION_CONTEXT(N'TenantRole') = N'Admin'
        AND EXISTS(
            SELECT 1
            FROM auth.Workspaces
            WHERE WorkspaceId = @WorkspaceId
                AND TenantDirectoryId = auth.fn_GetSessionTenantId()
        )
    );
GO;
CREATE SECURITY POLICY app.SpecificationsPolicy
ADD FILTER PREDICATE app.fn_Specifications_Filter(WorkspaceId) ON app.Specifications,
    ADD BLOCK PREDICATE app.fn_Specifications_Block(WorkspaceId) ON app.Specifications WITH (STATE = ON);
-- =============================================
GO;
CREATE FUNCTION app.fn_FileSchemas_Filter(@SpecificationId UNIQUEIDENTIFIER) RETURNS TABLE WITH SCHEMABINDING AS RETURN (
    WITH AssignedSpecification AS (
        SELECT WorkspaceId
        FROM app.Specifications
        WHERE SpecificationId = @SpecificationId
    )
    SELECT 1 AS Result
    FROM AssignedSpecification
        CROSS APPLY app.fn_Specifications_Filter(
            AssignedSpecification.WorkspaceId
        ) AS f
);
GO;
CREATE FUNCTION app.fn_FileSchemas_Block(@SpecificationId UNIQUEIDENTIFIER) RETURNS TABLE WITH SCHEMABINDING AS RETURN (
    WITH AssignedSpecification AS (
        SELECT WorkspaceId
        FROM app.Specifications
        WHERE SpecificationId = @SpecificationId
    )
    SELECT 1 AS Result
    FROM AssignedSpecification
        CROSS APPLY app.fn_Specifications_Block(
            AssignedSpecification.WorkspaceId
        ) AS f
);
GO;
CREATE SECURITY POLICY app.FileSchemasPolicy
ADD FILTER PREDICATE app.fn_FileSchemas_Filter(SpecificationId) ON app.FileSchemas,
    ADD BLOCK PREDICATE app.fn_FileSchemas_Block(SpecificationId) ON app.FileSchemas WITH (STATE = ON);
-- =============================================
GO;
CREATE FUNCTION app.fn_FileGroups_Filter(
    @SpecificationId UNIQUEIDENTIFIER,
    @CreatedBy UNIQUEIDENTIFIER,
    @ShareWithWorkspace BIT
) RETURNS TABLE WITH SCHEMABINDING AS RETURN
SELECT 1 AS Result
WHERE (
        @ShareWithWorkspace = 1
        AND EXISTS(
            SELECT 1
            FROM app.Specifications S
                INNER JOIN auth.UserWorkspaceAssignments UWA ON S.WorkspaceId = UWA.WorkspaceId
            WHERE S.SpecificationId = @SpecificationId
                AND UWA.UserObjectId = auth.fn_GetSessionUserId()
        )
    )
    OR (
        @ShareWithWorkspace = 0
        AND @CreatedBy = auth.fn_GetSessionUserId()
    )
    OR (
        SESSION_CONTEXT(N'TenantRole') = N'Admin'
        AND EXISTS(
            SELECT 1
            FROM app.Specifications S
                INNER JOIN auth.Workspaces W ON S.WorkspaceId = W.WorkspaceId
            WHERE S.SpecificationId = @SpecificationId
                AND W.TenantDirectoryId = auth.fn_GetSessionTenantId()
        )
    );
GO;
CREATE FUNCTION app.fn_FileGroups_Block(
    @SpecificationId UNIQUEIDENTIFIER,
    @CreatedBy UNIQUEIDENTIFIER,
    @ShareWithWorkspace BIT
) RETURNS TABLE WITH SCHEMABINDING AS RETURN
SELECT 1 AS Result
WHERE (
        @ShareWithWorkspace = 1
        AND EXISTS(
            SELECT 1
            FROM app.Specifications S
                INNER JOIN auth.UserWorkspaceAssignments UWA ON S.WorkspaceId = UWA.WorkspaceId
            WHERE S.SpecificationId = @SpecificationId
                AND UWA.UserObjectId = auth.fn_GetSessionUserId()
                AND (
                    UWA.RoleId = 1
                    OR UWA.RoleId = 2
                )
        )
    )
    OR (
        @ShareWithWorkspace = 0
        AND @CreatedBy = auth.fn_GetSessionUserId()
    )
    OR (
        SESSION_CONTEXT(N'TenantRole') = N'Admin'
        AND EXISTS(
            SELECT 1
            FROM app.Specifications S
                INNER JOIN auth.Workspaces W ON S.WorkspaceId = W.WorkspaceId
            WHERE S.SpecificationId = @SpecificationId
                AND W.TenantDirectoryId = auth.fn_GetSessionTenantId()
        )
    );
GO;
CREATE SECURITY POLICY app.FileGroupsPolicy
ADD FILTER PREDICATE app.fn_FileGroups_Filter(SpecificationId, CreatedBy, ShareWithWorkspace) ON app.FileGroups,
    ADD BLOCK PREDICATE app.fn_FileGroups_Block(SpecificationId, CreatedBy, ShareWithWorkspace) ON app.FileGroups WITH (STATE = ON);
-- =============================================
GO;
CREATE FUNCTION app.fn_Files_Filter(@FileGroupId UNIQUEIDENTIFIER) RETURNS TABLE WITH SCHEMABINDING AS RETURN (
    WITH AssignedFileGroup AS (
        SELECT SpecificationId,
            CreatedBy,
            ShareWithWorkspace
        FROM app.FileGroups
        WHERE FileGroupId = @FileGroupId
    )
    SELECT 1 AS Result
    FROM AssignedFileGroup
        CROSS APPLY app.fn_FileGroups_Filter(
            AssignedFileGroup.SpecificationId,
            AssignedFileGroup.CreatedBy,
            AssignedFileGroup.ShareWithWorkspace
        ) AS f
);
GO;
CREATE FUNCTION app.fn_Files_Block(@FileGroupId UNIQUEIDENTIFIER) RETURNS TABLE WITH SCHEMABINDING AS RETURN (
    WITH AssignedFileGroup AS (
        SELECT SpecificationId,
            CreatedBy,
            ShareWithWorkspace
        FROM app.FileGroups
        WHERE FileGroupId = @FileGroupId
    )
    SELECT 1 AS Result
    FROM AssignedFileGroup
        CROSS APPLY app.fn_FileGroups_Block(
            AssignedFileGroup.SpecificationId,
            AssignedFileGroup.CreatedBy,
            AssignedFileGroup.ShareWithWorkspace
        ) AS f
);
GO;
CREATE SECURITY POLICY app.FilesPolicy
ADD FILTER PREDICATE app.fn_Files_Filter(FileGroupId) ON app.Files,
    ADD BLOCK PREDICATE app.fn_Files_Block(FileGroupId) ON app.Files WITH (STATE = ON);
-- =============================================
GO;
CREATE FUNCTION app.fn_Models_Filter(@InputSpecificationId UNIQUEIDENTIFIER) RETURNS TABLE WITH SCHEMABINDING AS RETURN (
    WITH AssignedSpecification AS (
        SELECT WorkspaceId
        FROM app.Specifications
        WHERE SpecificationId = @InputSpecificationId
    )
    SELECT 1 AS Result
    FROM AssignedSpecification
        CROSS APPLY app.fn_Specifications_Filter(
            AssignedSpecification.WorkspaceId
        ) AS f
);
GO;
CREATE FUNCTION app.fn_Models_Block(@InputSpecificationId UNIQUEIDENTIFIER) RETURNS TABLE WITH SCHEMABINDING AS RETURN (
    WITH AssignedSpecification AS (
        SELECT WorkspaceId
        FROM app.Specifications
        WHERE SpecificationId = @InputSpecificationId
    )
    SELECT 1 AS Result
    FROM AssignedSpecification
        CROSS APPLY app.fn_Specifications_Block(
            AssignedSpecification.WorkspaceId
        ) AS f
);
GO;
CREATE SECURITY POLICY app.ModelsPolicy
ADD FILTER PREDICATE app.fn_Models_Filter(InputSpecificationId) ON app.Models,
    ADD BLOCK PREDICATE app.fn_Models_Block(InputSpecificationId) ON app.Models WITH (STATE = ON);
-- =============================================
GO;
CREATE FUNCTION app.fn_ModelRuns_Filter(@ModelId UNIQUEIDENTIFIER) RETURNS TABLE WITH SCHEMABINDING AS RETURN (
    WITH AssignedModel AS (
        SELECT InputSpecificationId
        FROM app.Models
        WHERE ModelId = @ModelId
    )
    SELECT 1 AS Result
    FROM AssignedModel
        CROSS APPLY app.fn_Models_Filter(
            AssignedModel.InputSpecificationId
        ) AS f
);
GO;
CREATE FUNCTION app.fn_ModelRuns_Block(@InputFileGroupId UNIQUEIDENTIFIER) RETURNS TABLE WITH SCHEMABINDING AS RETURN (
    WITH AssignedFileGroup AS (
        SELECT SpecificationId,
            CreatedBy,
            ShareWithWorkspace
        FROM app.FileGroups
        WHERE FileGroupId = @InputFileGroupId
    )
    SELECT 1 AS Result
    FROM AssignedFileGroup
        CROSS APPLY app.fn_FileGroups_Block(
            AssignedFileGroup.SpecificationId,
            AssignedFileGroup.CreatedBy,
            AssignedFileGroup.ShareWithWorkspace
        ) AS f
);
GO;
CREATE SECURITY POLICY app.ModelRunsPolicy
ADD FILTER PREDICATE app.fn_ModelRuns_Filter(ModelId) ON app.ModelRuns,
    ADD BLOCK PREDICATE app.fn_ModelRuns_Block(InputFileGroupId) ON app.ModelRuns WITH (STATE = ON);
-- =============================================
-- Users
-- =============================================
CREATE USER [iunvi-webapp]
FROM EXTERNAL PROVIDER;
-- Give the user the webapp role
ALTER ROLE WebApp
ADD MEMBER [iunvi-webapp];