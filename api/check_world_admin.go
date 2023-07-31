package api

import (
	"context"
	"database/sql"
	"fmt"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

func (server *Server) CheckWorldAdmin(ctx context.Context, worldId int32, needsSuperAdmin bool) error {

	err := server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
	if err == nil {
		return nil
	}

	authPayload, err := server.authorizeUserCookie(ctx)
	if err != nil {
		return unauthenticatedError(err)
	}

	isAdmin, err := server.store.IsWorldAdmin(ctx, db.IsWorldAdminParams{
		UserID:  authPayload.UserId,
		WorldID: worldId,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user is not an admin of this world")
		}
		return fmt.Errorf("failed to authorize world admin: %w", err)
	}

	if needsSuperAdmin {
		if isAdmin.SuperAdmin {
			return nil
		} else {
			return fmt.Errorf("SUPER ADMIN role required for this action")
		}
	}

	return nil
}
