package servicecore

import (
	"context"
	"database/sql"
	"fmt"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/token"
)

func (core *ServiceCore) CheckWorldPermissions(ctx context.Context, worldId int32, modulePermissions *ModulePermission) (*token.Payload, error) {
	module, err := core.Store.GetModule(ctx, db.GetModuleParams{
		WorldID: sql.NullInt32{
			Int32: worldId,
			Valid: true,
		},
	})

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("world does not exist")
		}
		return nil, fmt.Errorf("failed to get world : %w", err)
	}

	return core.CheckModuleIdPermissions(ctx, module.ID, modulePermissions)
}
