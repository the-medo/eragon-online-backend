package maps

import (
	"context"
	"github.com/the-medo/talebound-backend/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
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

	returnAllLayers := true
	_, err := server.CheckEntityTypePermissions(ctx, db.EntityTypeMap, req.GetMapId(), nil)
	if err != nil {
		returnAllLayers = false
	}

	layerRows, err := server.Store.GetMapLayers(ctx, req.GetMapId())
	if err != nil {
		return nil, err
	}

	rsp := &pb.GetMapLayersResponse{
		Layers: []*pb.ViewMapLayer{},
	}

	for _, layerRow := range layerRows {
		if returnAllLayers || layerRow.Enabled {
			rsp.Layers = append(rsp.Layers, converters.ConvertViewMapLayer(layerRow))
		}
	}

	return rsp, nil
}

func validateGetMapLayers(req *pb.GetMapLayersRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMapId(req.GetMapId()); err != nil {
		violations = append(violations, e.FieldViolation("map_id", err))
	}

	return violations
}
