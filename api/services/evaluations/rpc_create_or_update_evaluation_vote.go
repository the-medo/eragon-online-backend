package evaluations

import (
	"context"
	"database/sql"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *ServiceEvaluations) CreateOrUpdateEvaluationVote(ctx context.Context, req *pb.CreateOrUpdateEvaluationVoteRequest) (*pb.CreateOrUpdateEvaluationVoteResponse, error) {
	arg := db.GetEvaluationVoteByEvaluationIdUserIdAndVoterParams{
		EvaluationID: req.GetEvaluationId(),
		UserID:       req.GetUserId(),
		UserIDVoter:  req.GetUserIdVoter(),
	}

	var alreadyExists bool
	alreadyExists = true

	_, err := server.Store.GetEvaluationVoteByEvaluationIdUserIdAndVoter(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			alreadyExists = false
		} else {
			return nil, status.Errorf(codes.Internal, "failed to get evaluation vote: %v", err)
		}
	}

	rsp := &pb.CreateOrUpdateEvaluationVoteResponse{}

	if alreadyExists {

		evaluationVoteNew, err := server.Store.UpdateEvaluationVote(ctx, db.UpdateEvaluationVoteParams{
			EvaluationID: req.GetEvaluationId(),
			UserID:       req.GetUserId(),
			UserIDVoter:  req.GetUserIdVoter(),
			Value:        req.GetValue(),
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update evaluation vote: %v", err)
		}

		rsp.EvaluationVote = &pb.EvaluationVote{
			EvaluationId: evaluationVoteNew.EvaluationID,
			UserId:       evaluationVoteNew.UserID,
			UserIdVoter:  evaluationVoteNew.UserIDVoter,
			Value:        evaluationVoteNew.Value,
			CreatedAt:    timestamppb.New(evaluationVoteNew.CreatedAt),
		}
	} else {
		evaluationVoteNew, err := server.Store.CreateEvaluationVote(ctx, db.CreateEvaluationVoteParams{
			EvaluationID: req.GetEvaluationId(),
			UserID:       req.GetUserId(),
			UserIDVoter:  req.GetUserIdVoter(),
			Value:        req.GetValue(),
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create evaluation vote: %v", err)
		}

		rsp.EvaluationVote = &pb.EvaluationVote{
			EvaluationId: evaluationVoteNew.EvaluationID,
			UserId:       evaluationVoteNew.UserID,
			UserIdVoter:  evaluationVoteNew.UserIDVoter,
			Value:        evaluationVoteNew.Value,
			CreatedAt:    timestamppb.New(evaluationVoteNew.CreatedAt),
		}
	}

	return rsp, nil
}
