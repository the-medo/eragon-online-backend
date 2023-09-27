package locations

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

func (server *ServiceLocations) CreateWorldLocation(ctx context.Context, request *pb.CreateWorldLocationRequest) (*pb.ViewLocation, error) {
	violations := validateCreateWorldLocation(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	_, err := server.AuthorizeUserCookie(ctx)
	if err != nil {
		return nil, e.UnauthenticatedError(err)
	}

	argLocation := db.CreateLocationParams{
		Name: request.GetName(),
		Description: sql.NullString{
			String: request.GetDescription(),
			Valid:  true,
		},
		ThumbnailImageID: sql.NullInt32{
			Int32: request.GetThumbnailImageId(),
			Valid: true,
		},
	}

	location, err := server.Store.CreateLocation(ctx, argLocation)

	arg := db.CreateWorldLocationParams{
		WorldID:    request.GetWorldId(),
		LocationID: location.ID,
	}

	worldLocation, err := server.Store.CreateWorldLocation(ctx, arg)
	if err != nil {
		return nil, err
	}

	viewLocation, err := server.Store.GetLocationByID(ctx, worldLocation.LocationID)
	if err != nil {
		return nil, err
	}

	rsp := converters.ConvertViewLocation(viewLocation)

	return rsp, nil
}

func validateCreateWorldLocation(req *pb.CreateWorldLocationRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateWorldId(req.GetWorldId()); err != nil {
		violations = append(violations, e.FieldViolation("world_id", err))
	}

	if err := validator.ValidateLocationName(req.GetName()); err != nil {
		violations = append(violations, e.FieldViolation("name", err))
	}

	if req.Description != nil {
		if err := validator.ValidateLocationDescription(req.GetDescription()); err != nil {
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
