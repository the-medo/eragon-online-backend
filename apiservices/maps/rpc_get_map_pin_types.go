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

func (server *ServiceMaps) GetMapPinTypes(ctx context.Context, request *pb.GetMapPinTypesRequest) (*pb.GetMapPinTypesResponse, error) {
	violations := validateGetMapPinTypes(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	// currently works only for world map, not quest maps
	pinTypeRows, err := server.Store.GetMapPinTypesForMap(ctx, sql.NullInt32{Int32: request.GetMapId(), Valid: true})
	if err != nil {
		return nil, err
	}

	rsp := &pb.GetMapPinTypesResponse{
		PinTypes: make([]*pb.MapPinType, len(pinTypeRows)),
	}

	for i, pinTypeRow := range pinTypeRows {
		rsp.PinTypes[i] = converters.ConvertMapPinType(pinTypeRow)
	}

	return rsp, nil
}

func validateGetMapPinTypes(req *pb.GetMapPinTypesRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMapId(req.GetMapId()); err != nil {
		violations = append(violations, e.FieldViolation("map_id", err))
	}

	return violations
}
