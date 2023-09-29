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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *ServiceMaps) CreateWorldMap(ctx context.Context, request *pb.CreateWorldMapRequest) (*pb.CreateWorldMapResponse, error) {
	violations := validateCreateWorldMap(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	_, err := server.CheckWorldAdmin(ctx, request.GetWorldId(), false)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to add world tag: %v", err)
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
		Width:  request.GetWidth(),
		Height: request.GetHeight(),
		ThumbnailImageID: sql.NullInt32{
			Int32: request.GetThumbnailImageId(),
			Valid: request.ThumbnailImageId != nil,
		},
	}

	newMap, err := server.Store.CreateMap(ctx, argMap)
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

	rsp := &pb.CreateWorldMapResponse{
		Map:   converters.ConvertViewMap(viewMap),
		Layer: converters.ConvertViewMapLayer(viewMapLayer[0]),
	}

	return rsp, nil
}

func validateCreateWorldMap(req *pb.CreateWorldMapRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateWorldId(req.GetWorldId()); err != nil {
		violations = append(violations, e.FieldViolation("world_id", err))
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

	if err := validator.ValidateUniversalDimension(req.GetWidth()); err != nil {
		violations = append(violations, e.FieldViolation("width", err))
	}

	if err := validator.ValidateUniversalDimension(req.GetHeight()); err != nil {
		violations = append(violations, e.FieldViolation("height", err))
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
