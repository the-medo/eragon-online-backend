package locations

import (
	"context"
	"github.com/the-medo/talebound-backend/converters"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceLocations) GetLocations(ctx context.Context, req *pb.ModuleDefinition) (*pb.GetLocationsResponse, error) {
	violations := validateGetLocations(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	locationRows, err := server.Store.GetLocationsByModule(ctx, req.GetWorldId())
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

func validateGetLocations(req *pb.ModuleDefinition) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateLocationModule(req); err != nil {
		violations = append(violations, e.FieldViolation("modules", err))
	}

	return violations
}
