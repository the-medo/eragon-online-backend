package api

import (
	"context"
	"github.com/the-medo/talebound-backend/api/e"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) GetEvaluationVotesByUserIdAndVoter(ctx context.Context, req *pb.GetEvaluationVotesByUserIdAndVoterRequest) (*pb.GetEvaluationVotesByUserIdAndVoterResponse, error) {
	violations := validateGetEvaluationVotesByUserIdAndVoter(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	arg := db.GetEvaluationVotesByUserIdAndVoterParams{
		UserID:      req.GetUserId(),
		UserIDVoter: req.GetUserIdVoter(),
	}

	evaluationVotes, err := server.Store.GetEvaluationVotesByUserIdAndVoter(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get evaluation votes : %v", err)
	}

	rsp := &pb.GetEvaluationVotesByUserIdAndVoterResponse{
		EvaluationVote: make([]*pb.EvaluationVote, len(evaluationVotes)),
	}

	for i, evaluationVote := range evaluationVotes {
		rsp.EvaluationVote[i] = convertEvaluationVote(evaluationVote)
	}
	return rsp, nil
}

func validateGetEvaluationVotesByUserIdAndVoter(req *pb.GetEvaluationVotesByUserIdAndVoterRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUserId(req.GetUserId()); err != nil {
		violations = append(violations, e.FieldViolation("user_id", err))
	}
	if err := validator.ValidateUserId(req.GetUserIdVoter()); err != nil {
		violations = append(violations, e.FieldViolation("user_id_voter", err))
	}
	return violations
}
