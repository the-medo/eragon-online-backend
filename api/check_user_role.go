package api

import (
	"context"
	"database/sql"
	"fmt"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

func (server *Server) CheckUserRole(ctx context.Context, roleTypes []pb.RoleType) error {

	authPayload, err := server.authorizeUserCookie(ctx)
	if err != nil {
		return unauthenticatedError(err)
	}

	roleFound := false

	for _, roleType := range roleTypes {
		roles, err := server.store.HasUserRole(ctx, db.HasUserRoleParams{
			UserID: authPayload.UserId,
			Role:   roleType.String(),
		})

		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			return fmt.Errorf("failed to get user roles: %w", err)
		}

		if roles.UserID != authPayload.UserId {
			return fmt.Errorf("incorrect userId returned: %w", err)
		}

		roleFound = true
	}

	if !roleFound {
		return fmt.Errorf("user does not have any of the required roles")
	}

	return nil
}
