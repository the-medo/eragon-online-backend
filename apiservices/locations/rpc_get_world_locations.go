package locations

import (
	"context"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/api/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceLocations) GetWorldLocations(ctx context.Context, req *pb.GetWorldLocationsRequest) (*pb.GetWorldLocationResponse, error) {
	violations := validateGetWorldLocations(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	locationRows, err := server.Store.GetWorldLocations(ctx, req.GetWorldId())
	if err != nil {
		return nil, err
	}

	rsp := &pb.GetWorldLocationResponse{
		Locations: make([]*pb.ViewLocation, len(locationRows)),
	}

	for i, locationRow := range locationRows {
		rsp.Locations[i] = converters.ConvertViewLocation(locationRow)
	}

	return rsp, nil
}

func validateGetWorldLocations(req *pb.GetWorldLocationsRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateWorldId(req.GetWorldId()); err != nil {
		violations = append(violations, e.FieldViolation("world_id", err))
	}

	return violations
}
