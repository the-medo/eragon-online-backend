package users

import (
	"context"
	"github.com/the-medo/talebound-backend/api"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/api/e"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *ServiceUsers) GetUserPosts(ctx context.Context, req *pb.GetUserPostsRequest) (*pb.GetUserPostsResponse, error) {
	violations := validateGetUserPostsRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	limit, offset := api.GetDefaultQueryBoundaries(req.GetLimit(), req.GetOffset())

	authPayload, err := server.AuthorizeUserCookie(ctx)
	if err != nil {
		return nil, e.UnauthenticatedError(err)
	}

	if req.GetUserId() != authPayload.UserId {
		err := server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
		if err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "can not get list of posts - you are not creator or admin: %v", err)
		}
	}

	arg := db.GetPostsByUserIdParams{
		UserID:     req.GetUserId(),
		PageLimit:  limit,
		PageOffset: offset,
	}

	posts, err := server.Store.GetPostsByUserId(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get posts: %v", err)
	}

	rsp := &pb.GetUserPostsResponse{
		Posts: make([]*pb.ViewPost, len(posts)),
	}

	for i, post := range posts {
		rsp.Posts[i] = converters.ConvertViewPost(post)
	}

	return rsp, nil
}

// ================= VALIDATION =================

func validateGetUserPostsRequest(req *pb.GetUserPostsRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUserId(req.GetUserId()); err != nil {
		violations = append(violations, e.FieldViolation("user_id", err))
	}

	if req.Limit != nil {
		if err := validator.ValidateLimit(req.GetLimit()); err != nil {
			violations = append(violations, e.FieldViolation("limit", err))
		}
	}

	if req.Offset != nil {
		if err := validator.ValidateOffset(req.GetOffset()); err != nil {
			violations = append(violations, e.FieldViolation("offset", err))
		}
	}

	return violations
}
