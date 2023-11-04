package evaluations

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

func (server *ServiceEvaluations) DeleteEvaluationVote(ctx context.Context, req *pb.DeleteEvaluationVoteRequest) (*pb.DeleteEvaluationVoteResponse, error) {

	violations := validateDeleteEvaluationVote(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	arg := db.DeleteEvaluationVoteParams{
		EvaluationID: req.GetEvaluationId(),
		UserID:       req.GetUserId(),
		UserIDVoter:  req.GetUserIdVoter(),
	}

	err := server.Store.DeleteEvaluationVote(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get evaluation votes : %v", err)
	}

	rsp := &pb.DeleteEvaluationVoteResponse{
		Success: true,
		Message: "Evaluation vote deleted successfully",
	}

	return rsp, nil
}

func validateDeleteEvaluationVote(req *pb.DeleteEvaluationVoteRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateEvaluationId(req.GetEvaluationId()); err != nil {
		violations = append(violations, e.FieldViolation("evaluation_id", err))
	}

	if err := validator.ValidateUserId(req.GetUserId()); err != nil {
		violations = append(violations, e.FieldViolation("user_id", err))
	}

	if err := validator.ValidateUserId(req.GetUserIdVoter()); err != nil {
		violations = append(violations, e.FieldViolation("user_id_voter", err))
	}
	return violations
}
