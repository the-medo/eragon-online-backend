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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *ServiceLocations) CreateLocation(ctx context.Context, request *pb.CreateLocationRequest) (*pb.ViewLocation, error) {
	violations := validateCreateLocation(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	if request.GetModule().WorldId != nil {
		_, err := server.CheckWorldAdmin(ctx, request.GetModule().GetWorldId(), false)
		if err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "failed to create location: %v", err)
		}
	}

	if request.GetModule().QuestId != nil {
		return nil, status.Error(codes.Internal, "creating locations for quests is not implemented yet")
	}

	argLocation := db.CreateLocationParams{
		Name: request.GetName(),
		Description: sql.NullString{
			String: request.GetDescription(),
			Valid:  true,
		},
		ThumbnailImageID: sql.NullInt32{
			Int32: request.GetThumbnailImageId(),
			Valid: request.ThumbnailImageId != nil,
		},
	}

	location, err := server.Store.CreateLocation(ctx, argLocation)

	createdLocationId := int32(0)

	if request.GetModule().WorldId != nil {

		arg := db.CreateWorldLocationParams{
			WorldID:    request.GetModule().GetWorldId(),
			LocationID: location.ID,
		}

		worldLocation, err := server.Store.CreateWorldLocation(ctx, arg)
		if err != nil {
			return nil, err
		}

		createdLocationId = worldLocation.LocationID
	}

	viewLocation, err := server.Store.GetLocationByID(ctx, createdLocationId)
	if err != nil {
		return nil, err
	}

	rsp := converters.ConvertViewLocation(viewLocation)

	return rsp, nil
}

func validateCreateLocation(req *pb.CreateLocationRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateLocationModule(req.GetModule()); err != nil {
		violations = append(violations, e.FieldViolation("modules", err))
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
