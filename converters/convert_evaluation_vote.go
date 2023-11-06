package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertEvaluationVote(evaluationVote db.EvaluationVote) *pb.EvaluationVote {
	pbEvaluationVote := &pb.EvaluationVote{
		EvaluationId: evaluationVote.EvaluationID,
		UserId:       evaluationVote.UserID,
		UserIdVoter:  evaluationVote.UserIDVoter,
		Value:        evaluationVote.Value,
		CreatedAt:    timestamppb.New(evaluationVote.CreatedAt),
	}

	return pbEvaluationVote
}
