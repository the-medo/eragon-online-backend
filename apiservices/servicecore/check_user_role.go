package servicecore

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/the-medo/talebound-backend/api/e"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

func (core *ServiceCore) CheckUserRole(ctx context.Context, roleTypes []pb.RoleType) error {

	authPayload, err := core.AuthorizeUserCookie(ctx)
	if err != nil {
		return e.UnauthenticatedError(err)
	}

	roleFound := false

	for _, roleType := range roleTypes {
		roles, err := core.Store.HasUserRole(ctx, db.HasUserRoleParams{
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
