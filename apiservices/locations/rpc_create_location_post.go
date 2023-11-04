package locations

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/api/e"
	"github.com/the-medo/talebound-backend/apiservices/srv"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceLocations) CreateLocationPost(ctx context.Context, request *pb.CreateLocationPostRequest) (*pb.ViewLocation, error) {

	violations := validateCreateLocationPost(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	authPayload, err := server.CheckEntityTypePermissions(ctx, db.EntityTypeLocation, request.GetLocationId(), &srv.ModulePermission{
		NeedsEntityPermission: &[]db.EntityType{db.EntityTypeLocation, db.EntityTypePost},
	})

	location, err := server.Store.GetLocationByID(ctx, request.GetLocationId())
	if err != nil {
		return nil, err
	}

	newPost, err := server.Store.CreatePost(ctx, db.CreatePostParams{
		UserID:    authPayload.UserId,
		Title:     location.Name,
		IsDraft:   false,
		IsPrivate: false,
		Content:   "",
	})
	if err != nil {
		return nil, err
	}

	_, err = server.Store.UpdateLocation(ctx, db.UpdateLocationParams{
		ID: request.GetLocationId(),
		PostID: sql.NullInt32{
			Int32: newPost.ID,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	_, err = server.Store.CreateEntity(ctx, db.CreateEntityParams{
		Type:     db.EntityTypePost,
		ModuleID: location.ModuleID,
		PostID: sql.NullInt32{
			Int32: newPost.ID,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	rsp := converters.ConvertViewLocation(location)
	rsp.PostId = &newPost.ID
	rsp.PostTitle = &newPost.Title

	return rsp, nil
}

func validateCreateLocationPost(req *pb.CreateLocationPostRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateLocationId(req.GetLocationId()); err != nil {
		violations = append(violations, e.FieldViolation("location_id", err))
	}

	return violations
}
