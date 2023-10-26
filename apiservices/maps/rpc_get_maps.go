package maps

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/api/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceMaps) GetMaps(ctx context.Context, req *pb.GetMapsRequest) (*pb.GetMapsResponse, error) {
	violations := validateGetMaps(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	mapRows, err := server.Store.GetMaps(ctx, sql.NullInt32{
		Int32: req.GetPlacement().GetWorldId(),
		Valid: req.GetPlacement().GetWorldId() != 0,
	})

	if err != nil {
		return nil, err
	}

	rsp := &pb.GetMapsResponse{
		Maps: make([]*pb.ViewMap, len(mapRows)),
	}

	for i, mapRow := range mapRows {
		rsp.Maps[i] = converters.ConvertViewMap(mapRow)
	}

	return rsp, nil
}

func validateGetMaps(req *pb.GetMapsRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMapPlacement(req.GetPlacement()); err != nil {
		violations = append(violations, e.FieldViolation("placements", err))
	}

	return violations
}
