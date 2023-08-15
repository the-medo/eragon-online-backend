package api

import (
	"context"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) GetEvaluationVotesByUserId(ctx context.Context, req *pb.GetEvaluationVotesByUserIdRequest) (*pb.GetEvaluationVotesByUserIdResponse, error) {
	violations := validateGetEvaluationVotesByUserId(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	evaluationVotes, err := server.store.GetEvaluationVotesByUserId(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get evaluation votes : %v", err)
	}

	rsp := &pb.GetEvaluationVotesByUserIdResponse{
		EvaluationVote: make([]*pb.EvaluationVote, len(evaluationVotes)),
	}

	for i, evaluationVote := range evaluationVotes {
		rsp.EvaluationVote[i] = convertEvaluationVote(evaluationVote)
	}
	return rsp, nil
}

func validateGetEvaluationVotesByUserId(req *pb.GetEvaluationVotesByUserIdRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateUserId(req.GetUserId()); err != nil {
		violations = append(violations, FieldViolation("user_id", err))
	}
	return violations
}
