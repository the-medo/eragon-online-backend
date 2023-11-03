package api

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/api/e"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/util"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	violations := validateUpdateUser(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUserCookie(ctx)
	if err != nil {
		return nil, e.UnauthenticatedError(err)
	}

	if req.Id == 0 {
		req.Id = authPayload.UserId
	}

	if req.Id != authPayload.UserId {
		return nil, status.Errorf(codes.PermissionDenied, "you are not allowed to update this user")
	}

	if authPayload.UserId != req.GetId() {
		return nil, status.Errorf(codes.PermissionDenied, "you are not allowed to update this user")
	}

	arg := db.UpdateUserParams{
		ID: req.Id,
		Username: sql.NullString{
			String: req.GetUsername(),
			Valid:  req.Username != nil,
		},
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
		ImgID: sql.NullInt32{
			Int32: req.GetImgId(),
			Valid: req.ImgId != nil,
		},
		IntroductionPostID: sql.NullInt32{
			Int32: req.GetIntroductionPostId(),
			Valid: req.ImgId != nil,
		},
	}

	if req.Password != nil {
		hashedPassword, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
		}

		arg.HashedPassword = sql.NullString{
			String: hashedPassword,
			Valid:  true,
		}

		arg.PasswordChangedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}

	user, err := server.Store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	rsp := &pb.UpdateUserResponse{
		User: ConvertUserGetImage(server, ctx, user),
	}

	return rsp, nil
}

func validateUpdateUser(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUserId(req.GetId()); err != nil {
		violations = append(violations, e.FieldViolation("id", err))
	}

	if req.Username != nil {
		if err := validator.ValidateUsername(req.GetUsername()); err != nil {
			violations = append(violations, e.FieldViolation("username", err))
		}
	}

	if req.ImgId != nil {
		if err := validator.ValidateImageId(req.GetImgId()); err != nil {
			violations = append(violations, e.FieldViolation("img_id", err))
		}
	}

	if req.Password != nil {
		if err := validator.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, e.FieldViolation("password", err))
		}
	}

	if req.Email != nil {
		if err := validator.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, e.FieldViolation("email", err))
		}
	}

	if req.IntroductionPostId != nil {
		if err := validator.ValidatePostId(req.GetIntroductionPostId()); err != nil {
			violations = append(violations, e.FieldViolation("introduction_post_id", err))
		}
	}

	return violations
}
