package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

func ConvertViewModule(dbModule db.ViewModule) *pb.ViewModule {
	pbModule := &pb.ViewModule{
		ModuleId:   dbModule.ModuleID,
		ModuleType: ConvertModuleTypeToPB(dbModule.ModuleType),
		Tags:       dbModule.Tags,
		MenuId:     dbModule.MenuID,
	}

	if dbModule.ModuleWorldID.Valid {
		pbModule.ModuleWorldId = &dbModule.ModuleWorldID.Int32
	}
	if dbModule.ModuleWorldID.Valid {
		pbModule.ModuleWorldId = &dbModule.ModuleWorldID.Int32
	}
	if dbModule.ModuleQuestID.Valid {
		pbModule.ModuleQuestId = &dbModule.ModuleQuestID.Int32
	}
	if dbModule.ModuleCharacterID.Valid {
		pbModule.ModuleCharacterId = &dbModule.ModuleCharacterID.Int32
	}
	if dbModule.ModuleSystemID.Valid {
		pbModule.ModuleSystemId = &dbModule.ModuleSystemID.Int32
	}

	if dbModule.HeaderImgID.Valid {
		pbModule.HeaderImgId = dbModule.HeaderImgID.Int32
	}

	if dbModule.ThumbnailImgID.Valid {
		pbModule.ThumbnailImgId = dbModule.ThumbnailImgID.Int32
	}

	if dbModule.AvatarImgID.Valid {
		pbModule.AvatarImgId = dbModule.AvatarImgID.Int32
	}

	return pbModule
}
