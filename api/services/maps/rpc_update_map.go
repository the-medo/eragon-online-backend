package maps

import (
	"context"
	"database/sql"
	converters "github.com/the-medo/talebound-backend/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceMaps) UpdateMap(ctx context.Context, request *pb.UpdateMapRequest) (*pb.Map, error) {
	violations := validateUpdateMap(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	authPayload, err := server.CheckEntityTypePermissions(ctx, db.EntityTypeMap, request.GetMapId(), nil)
	if err != nil {
		return nil, err
	}

	argMap := db.UpdateMapParams{
		ID: request.GetMapId(),
		Title: sql.NullString{
			String: request.GetTitle(),
			Valid:  request.Title != nil,
		},
		Type: sql.NullString{
			String: request.GetType(),
			Valid:  request.Type != nil,
		},
		Description: sql.NullString{
			String: request.GetDescription(),
			Valid:  request.Description != nil,
		},
		ThumbnailImageID: sql.NullInt32{
			Int32: request.GetThumbnailImageId(),
			Valid: request.ThumbnailImageId != nil,
		},
		IsPrivate: sql.NullBool{
			Bool:  request.GetIsPrivate(),
			Valid: request.IsPrivate != nil,
		},
		LastUpdatedUserID: sql.NullInt32{
			Int32: authPayload.UserId,
			Valid: true,
		},
	}

	_, err = server.Store.UpdateMap(ctx, argMap)
	if err != nil {
		return nil, err
	}

	m, err := server.Store.GetMapById(ctx, request.GetMapId())
	if err != nil {
		return nil, err
	}

	rsp := converters.ConvertMap(m)

	return rsp, nil
}

func validateUpdateMap(req *pb.UpdateMapRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMapId(req.GetMapId()); err != nil {
		violations = append(violations, e.FieldViolation("map_id", err))
	}

	if req.Title != nil {
		if err := validator.ValidateUniversalName(req.GetTitle()); err != nil {
			violations = append(violations, e.FieldViolation("title", err))
		}
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

	if req.ThumbnailImageId != nil {
		if err := validator.ValidateImageId(req.GetThumbnailImageId()); err != nil {
			violations = append(violations, e.FieldViolation("thumbnail_image_id", err))
		}
	}

	return violations
}
