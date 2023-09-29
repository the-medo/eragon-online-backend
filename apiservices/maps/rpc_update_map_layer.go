package maps

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/api/e"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceMaps) UpdateMapLayer(ctx context.Context, request *pb.UpdateMapLayerRequest) (*pb.ViewMapLayer, error) {
	violations := validateUpdateMapLayer(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	err := server.CheckMapAccess(ctx, request.GetMapId(), false)
	if err != nil {
		return nil, err
	}

	argLayer := db.UpdateMapLayerParams{
		ID: request.GetLayerId(),
		Name: sql.NullString{
			String: request.GetName(),
			Valid:  request.Name != nil,
		},
		ImageID: sql.NullInt32{
			Int32: request.GetImageId(),
			Valid: request.ImageId != nil,
		},
		IsMain: sql.NullBool{
			Bool:  request.GetIsMain(),
			Valid: request.IsMain != nil,
		},
		Enabled: sql.NullBool{
			Bool:  request.GetEnabled(),
			Valid: request.Enabled != nil,
		},
		Sublayer: sql.NullBool{
			Bool:  request.GetSublayer(),
			Valid: request.Sublayer != nil,
		},
	}

	updatedLayer, err := server.Store.UpdateMapLayer(ctx, argLayer)
	if err != nil {
		return nil, err
	}

	viewMapLayer, err := server.Store.GetMapLayerByID(ctx, updatedLayer.ID)
	if err != nil {
		return nil, err
	}

	rsp := converters.ConvertViewMapLayer(viewMapLayer)

	return rsp, nil
}

func validateUpdateMapLayer(req *pb.UpdateMapLayerRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMapId(req.GetMapId()); err != nil {
		violations = append(violations, e.FieldViolation("map_id", err))
	}

	if err := validator.ValidateUniversalId(req.GetLayerId()); err != nil {
		violations = append(violations, e.FieldViolation("layer_id", err))
	}

	if req.Name != nil {
		if err := validator.ValidateUniversalName(req.GetName()); err != nil {
			violations = append(violations, e.FieldViolation("name", err))
		}
	}

	if req.ImageId != nil {
		if err := validator.ValidateImageId(req.GetImageId()); err != nil {
			violations = append(violations, e.FieldViolation("image_id", err))
		}
	}

	return violations
}
