package servicecore

import (
	"context"
	"database/sql"
	"fmt"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/token"
)

func (core *ServiceCore) CheckModuleIdPermissions(ctx context.Context, moduleId int32, modulePermissions *ModulePermission) (*token.Payload, error) {
	authPayload, err := core.AuthorizeUserCookie(ctx)
	if err != nil {
		return nil, e.UnauthenticatedError(err)
	}

	moduleAdmin, err := core.Store.GetModuleAdmin(ctx, db.GetModuleAdminParams{
		UserID:   authPayload.UserId,
		ModuleID: moduleId,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			err = core.CheckModuleAdmin(ctx, &moduleAdmin, modulePermissions)
			if err != nil {
				return nil, fmt.Errorf("user is not an admin of this module")
			} else {
				return authPayload, nil
			}
		}
		return nil, fmt.Errorf("failed to authorize module admin: %w", err)
	}

	return authPayload, core.CheckModuleAdmin(ctx, &moduleAdmin, modulePermissions)
}
