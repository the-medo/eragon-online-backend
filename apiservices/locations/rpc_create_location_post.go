package locations

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/api/e"
	"github.com/the-medo/talebound-backend/consts"
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

	authPayload, locationPlacement, err := server.CheckLocationAccess(ctx, request.GetLocationId(), false)
	if err != nil {
		return nil, err
	}

	location, err := server.Store.GetLocationByID(ctx, request.GetLocationId())
	if err != nil {
		return nil, err
	}

	newPost, err := server.Store.CreatePost(ctx, db.CreatePostParams{
		UserID:     authPayload.UserId,
		Title:      location.Name,
		PostTypeID: consts.PostTypeWorldDescription,
		IsDraft:    false,
		IsPrivate:  false,
		Content:    "",
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

	if locationPlacement.WorldId != nil && locationPlacement.GetWorldId() > 0 {
		_, err := server.Store.CreateWorldPost(ctx, db.CreateWorldPostParams{
			WorldID: locationPlacement.GetWorldId(),
			PostID:  newPost.ID,
		})
		if err != nil {
			return nil, err
		}
	}

	viewLocation, err := server.Store.GetLocationByID(ctx, request.GetLocationId())
	if err != nil {
		return nil, err
	}

	rsp := converters.ConvertViewLocation(viewLocation)

	return rsp, nil
}

func validateCreateLocationPost(req *pb.CreateLocationPostRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateLocationId(req.GetLocationId()); err != nil {
		violations = append(violations, e.FieldViolation("location_id", err))
	}

	return violations
}
