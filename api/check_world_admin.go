package api

import (
	"context"
	"database/sql"
	"fmt"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/token"
)

func (server *Server) CheckWorldAdmin(ctx context.Context, worldId int32, needsSuperAdmin bool) (*token.Payload, error) {
	authPayload, err := server.authorizeUserCookie(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	err = server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
	if err == nil {
		return authPayload, nil
	}

	isAdmin, err := server.store.IsWorldAdmin(ctx, db.IsWorldAdminParams{
		UserID:  authPayload.UserId,
		WorldID: worldId,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user is not an admin of this world")
		}
		return nil, fmt.Errorf("failed to authorize world admin: %w", err)
	}

	if needsSuperAdmin {
		if isAdmin.SuperAdmin {
			return authPayload, nil
		} else {
			return nil, fmt.Errorf("SUPER ADMIN role required for this action")
		}
	}

	return authPayload, nil
}
