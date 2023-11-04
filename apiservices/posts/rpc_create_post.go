package posts

import (
	"context"
	"github.com/the-medo/talebound-backend/api"
	"github.com/the-medo/talebound-backend/api/e"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *api.Server) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.Post, error) {
	violations := validateCreatePostRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUserCookie(ctx)
	if err != nil {
		return nil, e.UnauthenticatedError(err)
	}

	arg := db.CreatePostParams{
		UserID:     authPayload.UserId,
		PostTypeID: req.GetPostTypeId(),
		Title:      req.GetTitle(),
		Content:    req.GetContent(),
	}

	//in case we got a nil value for IsDraft or IsPrivate, we need to get the default value from the post type
	postTypeNeeded := false
	if req.IsDraft == nil || req.IsPrivate == nil {
		postTypeNeeded = true
	}

	if postTypeNeeded {
		postType, err := server.Store.GetPostTypeById(ctx, req.GetPostTypeId())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get post type: %v", err)
		}
		if req.IsDraft == nil {
			arg.IsDraft = postType.Draftable
		}
		if req.IsPrivate == nil {
			arg.IsPrivate = postType.Privatable
		}
	} else {
		arg.IsDraft = req.GetIsDraft()
		arg.IsPrivate = req.GetIsPrivate()
	}

	postResult, err := server.Store.CreatePost(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create post: %s", err)
	}

	postType, err := server.Store.GetPostTypeById(ctx, postResult.PostTypeID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get post type: %v", err)
	}

	viewPost, err := server.Store.GetPostById(ctx, postResult.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get post: %s", err)
	}

	rsp := api.convertPostAndPostType(viewPost, postType)

	return rsp, nil
}

func validateCreatePostRequest(req *pb.CreatePostRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidatePostTitle(req.GetTitle()); err != nil {
		violations = append(violations, e.FieldViolation("title", err))
	}

	if err := validator.ValidatePostContent(req.GetContent()); err != nil {
		violations = append(violations, e.FieldViolation("content", err))
	}

	if err := validator.ValidatePostTypeId(req.GetPostTypeId()); err != nil {
		violations = append(violations, e.FieldViolation("post_type_id", err))
	}

	return violations
}
