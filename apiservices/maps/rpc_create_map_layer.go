package maps

import (
	"context"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/api/e"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceMaps) CreateMapLayer(ctx context.Context, request *pb.CreateMapLayerRequest) (*pb.ViewMapLayer, error) {
	violations := validateCreateMapLayer(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	err := server.CheckMapAccess(ctx, request.GetMapId(), false)
	if err != nil {
		return nil, err
	}

	argLayer := db.CreateMapLayerParams{
		MapID:    request.GetMapId(),
		Name:     request.GetName(),
		ImageID:  request.GetImageId(),
		IsMain:   request.GetIsMain(),
		Enabled:  request.GetEnabled(),
		Sublayer: request.GetSublayer(),
	}

	newLayer, err := server.Store.CreateMapLayer(ctx, argLayer)
	if err != nil {
		return nil, err
	}

	viewMapLayer, err := server.Store.GetMapLayerByID(ctx, newLayer.ID)
	if err != nil {
		return nil, err
	}

	rsp := converters.ConvertViewMapLayer(viewMapLayer)

	return rsp, nil
}

func validateCreateMapLayer(req *pb.CreateMapLayerRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMapId(req.GetMapId()); err != nil {
		violations = append(violations, e.FieldViolation("map_id", err))
	}

	if err := validator.ValidateUniversalName(req.GetName()); err != nil {
		violations = append(violations, e.FieldViolation("name", err))
	}

	if err := validator.ValidateImageId(req.GetImageId()); err != nil {
		violations = append(violations, e.FieldViolation("image_id", err))
	}

	return violations
}
