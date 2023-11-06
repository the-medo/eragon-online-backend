package servicecore

import (
	"context"
	"database/sql"
	"fmt"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/token"
)

func (core *ServiceCore) CheckMenuPermissions(ctx context.Context, menuId int32, modulePermissions *ModulePermission) (*token.Payload, error) {
	authPayload, err := core.AuthorizeUserCookie(ctx)
	if err != nil {
		return nil, e.UnauthenticatedError(err)
	}

	err = core.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
	if err == nil {
		return authPayload, nil
	}

	moduleAdmin, err := core.Store.GetModuleAdminByMenuId(ctx, db.GetModuleAdminByMenuIdParams{
		UserID: authPayload.UserId,
		MenuID: sql.NullInt32{Int32: menuId, Valid: true},
	})

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("module admin by menu id not found")
		}
		return nil, fmt.Errorf("failed to get module admin by menu id: %w", err)
	}

	return authPayload, core.CheckModuleAdmin(ctx, &moduleAdmin, modulePermissions)
}
