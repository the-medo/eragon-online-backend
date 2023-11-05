package srv

import (
	"context"
	"database/sql"
	"fmt"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/token"
)

func (core *ServiceCore) CheckModuleTypePermissions(ctx context.Context, moduleType db.ModuleType, id int32, modulePermissions *ModulePermission) (*token.Payload, *db.Module, error) {
	module, err := core.Store.GetModule(ctx, db.GetModuleParams{
		WorldID: sql.NullInt32{
			Int32: id,
			Valid: moduleType == db.ModuleTypeWorld,
		},
		QuestID: sql.NullInt32{
			Int32: id,
			Valid: moduleType == db.ModuleTypeQuest,
		},
		SystemID: sql.NullInt32{
			Int32: id,
			Valid: moduleType == db.ModuleTypeSystem,
		},
		CharacterID: sql.NullInt32{
			Int32: id,
			Valid: moduleType == db.ModuleTypeCharacter,
		},
	})

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, fmt.Errorf("%v does not exist", moduleType)
		}
		return nil, nil, fmt.Errorf("failed to get world : %w", err)
	}

	authPayload, err := core.CheckModuleIdPermissions(ctx, module.ID, modulePermissions)
	return authPayload, &module, err
}
