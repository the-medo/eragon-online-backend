package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

var questStatusToPB = map[db.QuestStatus]pb.QuestStatus{
	db.QuestStatusUnknown:              pb.QuestStatus_UNKNOWN,
	db.QuestStatusNotStarted:           pb.QuestStatus_NOT_STARTED,
	db.QuestStatusInProgress:           pb.QuestStatus_IN_PROGRESS,
	db.QuestStatusFinishedCompleted:    pb.QuestStatus_FINISHED_COMPLETED,
	db.QuestStatusFinishedNotCompleted: pb.QuestStatus_FINISHED_NOT_COMPLETED,
}

var questStatusToDB = map[pb.QuestStatus]db.QuestStatus{
	pb.QuestStatus_UNKNOWN:                db.QuestStatusUnknown,
	pb.QuestStatus_NOT_STARTED:            db.QuestStatusNotStarted,
	pb.QuestStatus_IN_PROGRESS:            db.QuestStatusInProgress,
	pb.QuestStatus_FINISHED_COMPLETED:     db.QuestStatusFinishedCompleted,
	pb.QuestStatus_FINISHED_NOT_COMPLETED: db.QuestStatusFinishedNotCompleted,
}

func ConvertQuestStatusToPB(status db.QuestStatus) pb.QuestStatus {
	if val, ok := questStatusToPB[status]; ok {
		return val
	}
	return pb.QuestStatus_UNKNOWN
}

func ConvertQuestStatusToDB(status pb.QuestStatus) db.QuestStatus {
	if val, ok := questStatusToDB[status]; ok {
		return val
	}
	return db.QuestStatusUnknown
}
