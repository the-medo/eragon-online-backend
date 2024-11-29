package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertQuest(quest db.Quest) *pb.Quest {
	pbQuest := &pb.Quest{
		Id:               quest.ID,
		Name:             quest.Name,
		Public:           quest.Public,
		CreatedAt:        timestamppb.New(quest.CreatedAt),
		ShortDescription: quest.ShortDescription,
		Status:           ConvertQuestStatusToPB(quest.Status),
		CanJoin:          quest.CanJoin,
	}

	if quest.WorldID.Valid {
		pbQuest.WorldId = quest.WorldID.Int32
	}

	if quest.SystemID.Valid {
		pbQuest.SystemId = quest.SystemID.Int32
	}

	return pbQuest
}
