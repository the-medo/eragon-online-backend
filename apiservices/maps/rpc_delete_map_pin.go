package maps

import (
	"context"
	"github.com/the-medo/talebound-backend/api/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *ServiceMaps) DeleteMapPin(ctx context.Context, request *pb.DeleteMapPinRequest) (*emptypb.Empty, error) {
	violations := validateDeleteMapPin(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	err := server.CheckMapAccess(ctx, request.GetMapId(), false)
	if err != nil {
		return nil, err
	}

	err = server.Store.DeleteMapPin(ctx, request.GetPinId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func validateDeleteMapPin(req *pb.DeleteMapPinRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMapId(req.GetMapId()); err != nil {
		violations = append(violations, e.FieldViolation("map_id", err))
	}

	if err := validator.ValidateUniversalId(req.GetPinId()); err != nil {
		violations = append(violations, e.FieldViolation("pin_id", err))
	}

	return violations
}
