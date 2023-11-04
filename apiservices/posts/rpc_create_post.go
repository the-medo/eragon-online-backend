package posts

import (
	"context"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/api/e"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *ServicePosts) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.ViewPost, error) {
	violations := validateCreatePostRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	authPayload, err := server.AuthorizeUserCookie(ctx)
	if err != nil {
		return nil, e.UnauthenticatedError(err)
	}

	arg := db.CreatePostParams{
		UserID:    authPayload.UserId,
		Title:     req.GetTitle(),
		Content:   req.GetContent(),
		IsDraft:   req.GetIsDraft(),
		IsPrivate: req.GetIsPrivate(),
	}

	postResult, err := server.Store.CreatePost(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create post: %s", err)
	}

	viewPost, err := server.Store.GetPostById(ctx, postResult.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get post: %s", err)
	}

	rsp := converters.ConvertViewPost(viewPost)

	return rsp, nil
}

func validateCreatePostRequest(req *pb.CreatePostRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidatePostTitle(req.GetTitle()); err != nil {
		violations = append(violations, e.FieldViolation("title", err))
	}

	if err := validator.ValidatePostContent(req.GetContent()); err != nil {
		violations = append(violations, e.FieldViolation("content", err))
	}

	return violations
}
