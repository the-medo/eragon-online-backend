package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

func ConvertViewModule(dbModule db.ViewModule) *pb.ViewModule {
	pbModule := &pb.ViewModule{
		Id:                dbModule.ID,
		ModuleType:        ConvertModuleTypeToPB(dbModule.ModuleType),
		Tags:              dbModule.Tags,
		MenuId:            dbModule.MenuID,
		DescriptionPostId: dbModule.DescriptionPostID,
	}

	if dbModule.WorldID.Valid {
		pbModule.WorldId = &dbModule.WorldID.Int32
	}

	if dbModule.QuestID.Valid {
		pbModule.QuestId = &dbModule.QuestID.Int32
	}

	if dbModule.CharacterID.Valid {
		pbModule.CharacterId = &dbModule.CharacterID.Int32
	}

	if dbModule.SystemID.Valid {
		pbModule.SystemId = &dbModule.SystemID.Int32
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
