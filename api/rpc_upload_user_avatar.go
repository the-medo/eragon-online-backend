package api

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/the-medo/talebound-backend/api/e"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/*
UploadUserAvatar
Uploads a user avatar image to Cloudflare and inserts it into the DB
  - 1. filename: avatar-{userId}
  - 2. upload to cloudflare => get cloudflareId
  - 3. insert into DB `avatar-{userId}_{cloudflareId}`
  - 4. update user imgId in DB
*/
func (server *Server) UploadUserAvatar(ctx context.Context, request *pb.UploadUserAvatarRequest) (*pb.UploadUserAvatarResponse, error) {
	violations := validateUploadUserAvatarRequest(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUserCookie(ctx)
	if err != nil {
		return nil, e.UnauthenticatedError(err)
	}

	if authPayload.UserId != request.GetUserId() {
		return nil, status.Errorf(codes.PermissionDenied, "you are not allowed to update this user")
	}

	filename := fmt.Sprintf("avatar-%d", request.GetUserId())

	dbImg, err := server.UploadAndInsertToDb(ctx, request.GetData(), ImageTypeIdUserAvatar, filename, request.GetUserId())
	if err != nil {
		return nil, err
	}

	//update user imgID in DB
	_, err = server.Store.UpdateUser(ctx, db.UpdateUserParams{
		ID: request.GetUserId(),
		ImgID: sql.NullInt32{
			Int32: dbImg.ID,
			Valid: true,
		},
	})

	if err != nil {
		return nil, err
	}

	return &pb.UploadUserAvatarResponse{
		UserId: request.GetUserId(),
		Image:  ConvertImage(dbImg),
	}, nil
}

func validateUploadUserAvatarRequest(req *pb.UploadUserAvatarRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUserId(req.GetUserId()); err != nil {
		violations = append(violations, e.FieldViolation("user_id", err))
	}

	return violations
}
