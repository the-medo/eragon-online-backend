package users

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *ServiceUsers) UpdateUserIntroduction(ctx context.Context, req *pb.UpdateUserIntroductionRequest) (*pb.ViewPost, error) {
	violations := validateUpdateUserIntroductionRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	authPayload, err := server.AuthorizeUserCookie(ctx)
	if err != nil {
		return nil, e.UnauthenticatedError(err)
	}

	if req.GetUserId() != authPayload.UserId {
		err := server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
		if err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "unable to save changes - you are not creator or admin: %v", err)
		}
	}

	user, err := server.Store.GetUserById(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "unable to save changes - user not found: %s", err)
	}

	if !user.IntroductionPostID.Valid {
		//create new post
		createPostArg := db.CreatePostParams{
			UserID:    req.GetUserId(),
			Title:     "User introduction",
			Content:   req.GetContent(),
			IsDraft:   req.GetSaveAsDraft(),
			IsPrivate: false,
		}

		post, err := server.Store.CreatePost(ctx, createPostArg)
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
		_, err = server.Store.UpdateUser(ctx, updateUserArg)

		viewPost, err := server.Store.GetPostById(ctx, post.ID)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get post: %s", err)
		}

		return converters.ConvertViewPost(viewPost), nil
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
			IsDraft: sql.NullBool{
				Bool:  req.GetSaveAsDraft(),
				Valid: req.SaveAsDraft != nil,
			},
		}
		post, err := server.Store.UpdatePost(ctx, arg)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update post: %s", err)
		}

		viewPost, err := server.Store.GetPostById(ctx, post.ID)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get post: %s", err)
		}

		return converters.ConvertViewPost(viewPost), nil
	}
}

func validateUpdateUserIntroductionRequest(req *pb.UpdateUserIntroductionRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUserId(req.GetUserId()); err != nil {
		violations = append(violations, e.FieldViolation("user_id", err))
	}

	if err := validator.ValidatePostContent(req.GetContent()); err != nil {
		violations = append(violations, e.FieldViolation("content", err))
	}

	return violations
}
