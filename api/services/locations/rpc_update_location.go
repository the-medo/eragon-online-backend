package locations

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/api/servicecore"
	converters "github.com/the-medo/talebound-backend/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceLocations) UpdateLocation(ctx context.Context, request *pb.UpdateLocationRequest) (*pb.Location, error) {
	violations := validateUpdateLocation(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	_, err := server.CheckEntityTypePermissions(ctx, db.EntityTypeLocation, request.GetLocationId(), &servicecore.ModulePermission{
		NeedsEntityPermission: &[]db.EntityType{db.EntityTypeLocation},
	})
	if err != nil {
		return nil, err
	}

	argLocation := db.UpdateLocationParams{
		ID: request.GetLocationId(),
		Name: sql.NullString{
			String: request.GetName(),
			Valid:  request.Name != nil,
		},
		Description: sql.NullString{
			String: request.GetDescription(),
			Valid:  request.Description != nil,
		},
		PostID: sql.NullInt32{
			Int32: request.GetPostId(),
			Valid: request.PostId != nil,
		},
		ThumbnailImageID: sql.NullInt32{
			Int32: request.GetThumbnailImageId(),
			Valid: request.ThumbnailImageId != nil,
		},
	}

	location, err := server.Store.UpdateLocation(ctx, argLocation)
	if err != nil {
		return nil, err
	}

	rsp := converters.ConvertLocation(location)

	return rsp, nil
}

func validateUpdateLocation(req *pb.UpdateLocationRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateLocationId(req.GetLocationId()); err != nil {
		violations = append(violations, e.FieldViolation("location_id", err))
	}

	if req.Name != nil {
		if err := validator.ValidateLocationName(req.GetName()); err != nil {
			violations = append(violations, e.FieldViolation("name", err))
		}
	}

	if req.Description != nil {
		if err := validator.ValidateLocationDescription(req.GetDescription()); err != nil {
			violations = append(violations, e.FieldViolation("description", err))
		}
	}

	if req.PostId != nil {
		if err := validator.ValidatePostId(req.GetPostId()); err != nil {
			violations = append(violations, e.FieldViolation("post_id", err))
		}
	}

	if req.ThumbnailImageId != nil {
		if err := validator.ValidateImageId(req.GetThumbnailImageId()); err != nil {
			violations = append(violations, e.FieldViolation("thumbnail_image_id", err))
		}
	}

	return violations
}
