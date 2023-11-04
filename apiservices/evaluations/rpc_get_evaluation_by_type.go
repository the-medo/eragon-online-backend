package evaluations

import (
	"context"
	"github.com/the-medo/talebound-backend/api"
	"github.com/the-medo/talebound-backend/api/e"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/util"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *api.Server) GetEvaluationsByType(ctx context.Context, req *pb.GetEvaluationsByTypeRequest) (*pb.GetEvaluationsByTypeResponse, error) {

	violations := validateGetEvaluationsByType(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	evaluations, err := server.Store.GetEvaluationsByType(ctx, db.EvaluationType(req.GetType()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get evaluation list: %v", err)
	}

	rsp := &pb.GetEvaluationsByTypeResponse{
		Evaluation: make([]*pb.Evaluation, len(evaluations)),
	}

	for i, evaluation := range evaluations {
		evaluationType, err := util.StringToEvaluationType(string(evaluation.EvaluationType))
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to convert evaluation type: %v", err)
		}

		rsp.Evaluation[i] = &pb.Evaluation{
			Id:          evaluation.ID,
			Name:        evaluation.Name,
			Description: evaluation.Description,
			Type:        evaluationType,
		}
	}

	return rsp, nil
}

func validateGetEvaluationsByType(req *pb.GetEvaluationsByTypeRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	return violations
}
