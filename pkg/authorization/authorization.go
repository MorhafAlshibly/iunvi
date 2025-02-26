package authorization

import (
	"context"

	"github.com/MorhafAlshibly/iunvi/gen/api"
	"github.com/MorhafAlshibly/iunvi/internal/tenantManagement/model"
	"github.com/MorhafAlshibly/iunvi/pkg/conversion"
	"github.com/MorhafAlshibly/iunvi/pkg/middleware"
	mssql "github.com/microsoft/go-mssqldb"
)

type Authorization struct {
	workspaceId   mssql.UniqueIdentifier
	workspaceRole api.WorkspaceRole
}

func WithWorkspaceID(workspaceId mssql.UniqueIdentifier) func(*Authorization) {
	return func(input *Authorization) {
		input.workspaceId = workspaceId
	}
}

func WithWorkspaceRole(workspaceRole api.WorkspaceRole) func(*Authorization) {
	return func(input *Authorization) {
		input.workspaceRole = workspaceRole
	}
}

func NewAuthorization(options ...func(*Authorization)) *Authorization {
	authorization := &Authorization{}
	for _, option := range options {
		option(authorization)
	}
	return authorization
}

func (d *Authorization) IsAuthorized(ctx context.Context) (bool, error) {
	tenantRole := ctx.Value("TenantRole").(string)
	if tenantRole == "Admin" {
		return true, nil
	}
	database := model.New(middleware.GetTx(ctx))
	userObjectId := ctx.Value("UserObjectId").(string)
	userObjectIdBytes, err := conversion.StringToUniqueIdentifier(userObjectId)
	if err != nil {
		return false, err
	}
	check, err := database.AuthorizationCheck(ctx, model.AuthorizationCheckParams{
		WorkspaceId:  d.workspaceId,
		UserObjectId: userObjectIdBytes,
		RoleName:     d.workspaceRole.String(),
	})
	if err != nil {
		return false, err
	}
	if check > 0 {
		return true, nil
	}
	return false, nil
}
