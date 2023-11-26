package locations

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

func (server *ServiceLocations) GetLocationById(ctx context.Context, req *pb.GetLocationByIdRequest) (*pb.Location, error) {
	violations := validateGetLocationByIdRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	m, err := server.Store.GetLocationById(ctx, req.GetLocationId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get location: %v", err)
	}

	return converters.ConvertLocation(m), nil
}

func validateGetLocationByIdRequest(req *pb.GetLocationByIdRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUniversalId(req.GetLocationId()); err != nil {
		violations = append(violations, e.FieldViolation("location_id", err))
	}
	return violations
}
