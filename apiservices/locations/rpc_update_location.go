package locations

import (
	"context"
	"database/sql"
	converters "github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/api/e"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *ServiceLocations) UpdateLocation(ctx context.Context, request *pb.UpdateLocationRequest) (*pb.ViewLocation, error) {
	violations := validateUpdateLocation(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	assignments, err := server.Store.GetLocationAssignments(ctx, request.GetLocationId())
	if assignments.WorldID > 0 {
		_, err = server.CheckWorldAdmin(ctx, assignments.WorldID, false)
		if err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "failed to update location: %v", err)
		}
	}

	_, err = server.AuthorizeUserCookie(ctx)
	if err != nil {
		return nil, e.UnauthenticatedError(err)
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

	_, err = server.Store.UpdateLocation(ctx, argLocation)
	if err != nil {
		return nil, err
	}

	location, err := server.Store.GetLocationByID(ctx, request.GetLocationId())
	if err != nil {
		return nil, err
	}

	rsp := converters.ConvertViewLocation(location)

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
