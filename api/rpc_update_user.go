package api

import (
	"context"
	"database/sql"
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
	violations := validateUpdateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUserCookie(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
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

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	rsp := &pb.UpdateUserResponse{
		User: convertUserGetImage(server, ctx, user),
	}

	return rsp, nil
}

func (server *Server) UpdateUserIntroduction(ctx context.Context, req *pb.UpdateUserIntroductionRequest) (*pb.Post, error) {
	violations := validateUpdateUserIntroductionRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUserCookie(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if req.GetUserId() != authPayload.UserId {
		err := server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
		if err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "can not update this post - you are not creator or admin: %v", err)
		}
	}

	user, err := server.store.GetUserById(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
	}

	if !user.IntroductionPostID.Valid {
		//create new post
		createPostArg := db.CreatePostParams{
			UserID:     req.GetUserId(),
			Title:      "User introduction",
			PostTypeID: 800,
			Content:    req.GetContent(),
		}

		post, err := server.store.CreatePost(ctx, createPostArg)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create post: %s", err)
		}

		//if it is post (so first time introduction is created), we put it into users table
		updateUserArg := db.UpdateUserParams{
			ID: req.UserId,
			IntroductionPostID: sql.NullInt32{
				Int32: post.ID,
				Valid: true,
			},
		}
		_, err = server.store.UpdateUser(ctx, updateUserArg)

		return convertPost(post), nil
	} else {
		//update existing post
		arg := db.UpdatePostParams{
			PostID: user.IntroductionPostID.Int32,
			LastUpdatedUserID: sql.NullInt32{
				Int32: authPayload.UserId,
				Valid: true,
			},
			Content: sql.NullString{
				String: req.GetContent(),
				Valid:  true,
			},
		}
		post, err := server.store.UpdatePost(ctx, arg)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update post: %s", err)
		}

		return convertPost(post), nil
	}
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUserId(req.GetId()); err != nil {
		violations = append(violations, FieldViolation("id", err))
	}

	if req.Username != nil {
		if err := validator.ValidateUsername(req.GetUsername()); err != nil {
			violations = append(violations, FieldViolation("username", err))
		}
	}

	if req.ImgId != nil {
		if err := validator.ValidateImgId(req.GetImgId()); err != nil {
			violations = append(violations, FieldViolation("img_id", err))
		}
	}

	if req.Password != nil {
		if err := validator.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, FieldViolation("password", err))
		}
	}

	if req.Email != nil {
		if err := validator.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, FieldViolation("email", err))
		}
	}

	if req.IntroductionPostId != nil {
		if err := validator.ValidatePostId(req.GetIntroductionPostId()); err != nil {
			violations = append(violations, FieldViolation("introduction_post_id", err))
		}
	}

	return violations
}

func validateUpdateUserIntroductionRequest(req *pb.UpdateUserIntroductionRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUserId(req.GetUserId()); err != nil {
		violations = append(violations, FieldViolation("user_id", err))
	}

	return violations
}
