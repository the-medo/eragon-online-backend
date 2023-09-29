package maps

import (
	"context"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/api/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceMaps) GetWorldMaps(ctx context.Context, req *pb.GetWorldMapRequest) (*pb.GetWorldMapResponse, error) {
	violations := validateGetWorldMaps(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	mapRows, err := server.Store.GetWorldMaps(ctx, req.GetWorldId())
	if err != nil {
		return nil, err
	}

	rsp := &pb.GetWorldMapResponse{
		Maps: make([]*pb.ViewMap, len(mapRows)),
	}

	for i, mapRow := range mapRows {
		rsp.Maps[i] = converters.ConvertViewMap(mapRow)
	}

	return rsp, nil
}

func validateGetWorldMaps(req *pb.GetWorldMapRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateWorldId(req.GetWorldId()); err != nil {
		violations = append(violations, e.FieldViolation("world_id", err))
	}

	return violations
}
