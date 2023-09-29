package maps

import (
	"context"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/api/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceMaps) GetMapPins(ctx context.Context, req *pb.GetMapPinsRequest) (*pb.GetMapPinsResponse, error) {
	violations := validateGetMapPins(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	pinRows, err := server.Store.GetMapPins(ctx, req.GetMapId())
	if err != nil {
		return nil, err
	}

	rsp := &pb.GetMapPinsResponse{
		Pins: make([]*pb.ViewMapPin, len(pinRows)),
	}

	for i, pinRow := range pinRows {
		rsp.Pins[i] = converters.ConvertViewMapPin(pinRow)
	}

	return rsp, nil
}

func validateGetMapPins(req *pb.GetMapPinsRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMapId(req.GetMapId()); err != nil {
		violations = append(violations, e.FieldViolation("map_id", err))
	}

	return violations
}
