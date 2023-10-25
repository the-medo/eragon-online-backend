package locations

import (
	"context"
	"github.com/the-medo/talebound-backend/api/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *ServiceLocations) DeleteLocation(ctx context.Context, request *pb.DeleteLocationRequest) (*emptypb.Empty, error) {
	violations := validateDeleteLocation(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	_, _, err := server.CheckLocationAccess(ctx, request.GetLocationId(), false)
	if err != nil {
		return nil, err
	}

	err = server.Store.DeleteLocation(ctx, request.GetLocationId())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func validateDeleteLocation(req *pb.DeleteLocationRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateLocationId(req.GetLocationId()); err != nil {
		violations = append(violations, e.FieldViolation("location_id", err))
	}
	return violations
}
