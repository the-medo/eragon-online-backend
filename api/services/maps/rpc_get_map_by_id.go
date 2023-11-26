package maps

import (
	"context"
	"github.com/the-medo/talebound-backend/converters"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *ServiceMaps) GetMapById(ctx context.Context, req *pb.GetMapByIdRequest) (*pb.Map, error) {
	violations := validateGetMapByIdRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	m, err := server.Store.GetMapById(ctx, req.GetMapId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get map: %v", err)
	}

	return converters.ConvertMap(m), nil
}

func validateGetMapByIdRequest(req *pb.GetMapByIdRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUniversalId(req.GetMapId()); err != nil {
		violations = append(violations, e.FieldViolation("map_id", err))
	}
	return violations
}
