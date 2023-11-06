package servicecore

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/the-medo/talebound-backend/api/e"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/token"
)

func (core *ServiceCore) CheckEntityPermissions(ctx context.Context, entityId int32, modulePermissions *ModulePermission) (*token.Payload, error) {
	authPayload, err := core.AuthorizeUserCookie(ctx)
	if err != nil {
		return nil, e.UnauthenticatedError(err)
	}

	moduleAdmin, err := core.Store.GetEntityModuleAdmin(ctx, db.GetEntityModuleAdminParams{
		UserID:   authPayload.UserId,
		EntityID: entityId,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user is not an admin of entity module")
		}
		return nil, fmt.Errorf("failed to authorize module admin for entity: %w", err)
	}

	return authPayload, core.CheckModuleAdmin(ctx, &moduleAdmin, modulePermissions)
}
