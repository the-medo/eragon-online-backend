package locations

import (
	"context"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/api/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceLocations) GetLocations(ctx context.Context, req *pb.LocationPlacement) (*pb.GetLocationsResponse, error) {
	violations := validateGetLocations(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	locationRows, err := server.Store.GetLocationsForPlacement(ctx, req.GetWorldId())
	if err != nil {
		return nil, err
	}

	rsp := &pb.GetLocationsResponse{
		Locations: make([]*pb.ViewLocation, len(locationRows)),
	}

	for i, locationRow := range locationRows {
		rsp.Locations[i] = converters.ConvertViewLocation(locationRow)
	}

	return rsp, nil
}

func validateGetLocations(req *pb.LocationPlacement) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateLocationPlacement(req); err != nil {
		violations = append(violations, e.FieldViolation("placement", err))
	}

	return violations
}
