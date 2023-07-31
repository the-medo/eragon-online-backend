package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertWorldActivity(dbWorldActivity db.WorldActivity) *pb.WorldActivity {
	pbWorldActivity := &pb.WorldActivity{
		WorldId:               dbWorldActivity.WorldID,
		Date:                  timestamppb.New(dbWorldActivity.Date),
		ActivityPostCount:     dbWorldActivity.PostCount,
		ActivityQuestCount:    dbWorldActivity.QuestCount,
		ActivityResourceCount: dbWorldActivity.ResourceCount,
	}

	return pbWorldActivity
}
