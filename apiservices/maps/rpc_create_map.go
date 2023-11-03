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

func (server *ServiceMaps) CreateMap(ctx context.Context, request *pb.CreateMapRequest) (*pb.CreateMapResponse, error) {
	violations := validateCreateMap(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	_, err := server.CheckModuleIdPermissions(ctx, request.GetModuleId(), nil)
	if err != nil {
		return nil, err
	}

	img, err := server.Store.GetImageById(ctx, request.GetLayerImageId())
	if err != nil {
		return nil, err
	}

	argMap := db.CreateMapParams{
		Name: request.GetName(),
		Type: sql.NullString{
			String: request.GetType(),
			Valid:  request.Type != nil,
		},
		Description: sql.NullString{
			String: request.GetDescription(),
			Valid:  request.Description != nil,
		},
		Width:  img.Width,
		Height: img.Height,
		ThumbnailImageID: sql.NullInt32{
			Int32: request.GetThumbnailImageId(),
			Valid: request.ThumbnailImageId != nil,
		},
	}

	newMap, err := server.Store.CreateMap(ctx, argMap)
	if err != nil {
		return nil, err
	}

	_, err = server.Store.CreateEntity(ctx, db.CreateEntityParams{
		Type:     db.EntityTypeMap,
		ModuleID: request.GetModuleId(),
		MapID: sql.NullInt32{
			Int32: newMap.ID,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	argMapLayer := db.CreateMapLayerParams{
		Name:     request.GetName(),
		MapID:    newMap.ID,
		ImageID:  request.GetLayerImageId(),
		IsMain:   true,
		Enabled:  true,
		Sublayer: false,
	}

	// Assuming a function to create the main layer for the map based on LayerImageID
	_, err = server.Store.CreateMapLayer(ctx, argMapLayer)
	if err != nil {
		return nil, err
	}

	viewMap, err := server.Store.GetMapByID(ctx, newMap.ID)
	if err != nil {
		return nil, err
	}

	viewMapLayer, err := server.Store.GetMapLayers(ctx, newMap.ID)
	if err != nil {
		return nil, err
	}

	rsp := &pb.CreateMapResponse{
		Map:   converters.ConvertViewMap(viewMap),
		Layer: converters.ConvertViewMapLayer(viewMapLayer[0]),
	}

	return rsp, nil
}

func validateCreateMap(req *pb.CreateMapRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateUniversalId(req.GetModuleId()); err != nil {
		violations = append(violations, e.FieldViolation("module_id", err))
	}

	if err := validator.ValidateUniversalName(req.GetName()); err != nil {
		violations = append(violations, e.FieldViolation("name", err))
	}

	if req.Type != nil {
		if err := validator.ValidateUniversalName(req.GetType()); err != nil {
			violations = append(violations, e.FieldViolation("type", err))
		}
	}

	if req.Description != nil {
		if err := validator.ValidateUniversalDescription(req.GetDescription()); err != nil {
			violations = append(violations, e.FieldViolation("description", err))
		}
	}

	if err := validator.ValidateImageId(req.GetLayerImageId()); err != nil {
		violations = append(violations, e.FieldViolation("layer_image_id", err))
	}

	if req.ThumbnailImageId != nil {
		if err := validator.ValidateImageId(req.GetThumbnailImageId()); err != nil {
			violations = append(violations, e.FieldViolation("thumbnail_image_id", err))
		}
	}

	return violations
}
