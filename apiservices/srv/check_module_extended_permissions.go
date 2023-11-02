package srv

import (
	"context"
	"fmt"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/token"
)

func (core *ServiceCore) CheckModuleExtendedPermissions(ctx context.Context, module *pb.Module, modulePermissions *ModulePermission) (*token.Payload, error) {
	if module == nil {
		return nil, fmt.Errorf("module is nil")
	}

	if module.WorldId != nil && module.GetWorldId() > 0 {
		return core.CheckModuleTypePermissions(ctx, db.ModuleTypeWorld, module.GetWorldId(), modulePermissions)
	}

	if module.QuestId != nil && module.GetQuestId() > 0 {
		return core.CheckModuleTypePermissions(ctx, db.ModuleTypeQuest, module.GetQuestId(), modulePermissions)
	}

	if module.SystemId != nil && module.GetSystemId() > 0 {
		return core.CheckModuleTypePermissions(ctx, db.ModuleTypeSystem, module.GetSystemId(), modulePermissions)
	}

	if module.CharacterId != nil && module.GetCharacterId() > 0 {
		return core.CheckModuleTypePermissions(ctx, db.ModuleTypeCharacter, module.GetCharacterId(), modulePermissions)
	}

	return nil, fmt.Errorf("module does not have a valid id")
}
