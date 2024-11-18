package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertViewQuest(quest db.ViewQuest) *pb.ViewQuest {
	pbViewQuest := &pb.ViewQuest{
		Id:               quest.ID,
		ModuleId:         quest.ModuleID,
		Name:             quest.Name,
		Public:           quest.Public,
		CreatedAt:        timestamppb.New(quest.CreatedAt),
		ShortDescription: quest.ShortDescription,
		Tags:             quest.Tags,
		MenuId:           quest.MenuID,
		CanJoin:          quest.CanJoin,
		Status:           ConvertQuestStatusToPB(quest.Status),
	}

	if quest.WorldID.Valid {
		pbViewQuest.WorldId = quest.WorldID.Int32
	}

	if quest.SystemID.Valid {
		pbViewQuest.SystemId = quest.SystemID.Int32
	}

	if quest.HeaderImgID.Valid {
		pbViewQuest.HeaderImgId = quest.HeaderImgID.Int32
	}

	if quest.ThumbnailImgID.Valid {
		pbViewQuest.ThumbnailImgId = quest.ThumbnailImgID.Int32
	}

	if quest.AvatarImgID.Valid {
		pbViewQuest.AvatarImgId = quest.AvatarImgID.Int32
	}

	return pbViewQuest
}
