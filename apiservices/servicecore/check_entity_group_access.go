package servicecore

import (
	"context"
	"github.com/the-medo/talebound-backend/api/e"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (core *ServiceCore) CheckEntityGroupAccess(ctx context.Context, entityGroupId int32, modulePermissions *ModulePermission) error {

	menuId, err := core.Store.GetMenuIdOfEntityGroup(ctx, entityGroupId)
	_, err = core.CheckMenuPermissions(ctx, menuId, modulePermissions)
	if err != nil {
		return status.Errorf(codes.PermissionDenied, "failed to get entity group access - not menu admin: %v", err)
	}

	_, err = core.AuthorizeUserCookie(ctx)
	if err != nil {
		return e.UnauthenticatedError(err)
	}

	return nil
}
