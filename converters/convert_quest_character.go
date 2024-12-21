package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertQuestCharacter(dbQuestCharacter db.QuestCharacter) *pb.QuestCharacter {
	pbQuestCharacter := &pb.QuestCharacter{
		QuestId:            dbQuestCharacter.QuestID,
		CharacterId:        dbQuestCharacter.CharacterID,
		CreatedAt:          timestamppb.New(dbQuestCharacter.CreatedAt),
		Approved:           dbQuestCharacter.Approved,
		MotivationalLetter: dbQuestCharacter.MotivationalLetter,
	}

	return pbQuestCharacter
}
