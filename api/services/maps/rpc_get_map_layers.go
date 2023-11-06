package maps

import (
	"context"
	"github.com/the-medo/talebound-backend/converters"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceMaps) GetMapLayers(ctx context.Context, req *pb.GetMapLayersRequest) (*pb.GetMapLayersResponse, error) {
	violations := validateGetMapLayers(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	layerRows, err := server.Store.GetMapLayers(ctx, req.GetMapId())
	if err != nil {
		return nil, err
	}

	rsp := &pb.GetMapLayersResponse{
		Layers: make([]*pb.ViewMapLayer, len(layerRows)),
	}

	for i, layerRow := range layerRows {
		rsp.Layers[i] = converters.ConvertViewMapLayer(layerRow)
	}

	return rsp, nil
}

func validateGetMapLayers(req *pb.GetMapLayersRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMapId(req.GetMapId()); err != nil {
		violations = append(violations, e.FieldViolation("map_id", err))
	}

	return violations
}
