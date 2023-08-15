package api

import (
	"context"
	"database/sql"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) GetUserPosts(ctx context.Context, req *pb.GetUserPostsRequest) (*pb.GetUserPostsResponse, error) {
	violations := validateGetUserPostsRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	limit, offset := GetDefaultQueryBoundaries(req.GetLimit(), req.GetOffset())

	authPayload, err := server.authorizeUserCookie(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if req.GetUserId() != authPayload.UserId {
		err := server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
		if err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "can not get list of posts - you are not creator or admin: %v", err)
		}
	}

	arg := db.GetPostsByUserIdParams{
		UserID: req.GetUserId(),
		PostTypeID: sql.NullInt32{
			Int32: req.GetPostTypeId(),
			Valid: req.PostTypeId != nil,
		},
		PageLimit:  limit,
		PageOffset: offset,
	}

	posts, err := server.store.GetPostsByUserId(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get posts: %v", err)
	}

	rsp := &pb.GetUserPostsResponse{
		Posts: make([]*pb.Post, len(posts)),
	}

	for i, post := range posts {
		rsp.Posts[i] = convertViewPost(post)
	}

	return rsp, nil
}

// ================= VALIDATION =================

func validateGetUserPostsRequest(req *pb.GetUserPostsRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUserId(req.GetUserId()); err != nil {
		violations = append(violations, FieldViolation("user_id", err))
	}

	if req.PostTypeId != nil {
		if err := validator.ValidatePostTypeId(req.GetPostTypeId()); err != nil {
			violations = append(violations, FieldViolation("post_type_id", err))
		}
	}

	if req.Limit != nil {
		if err := validator.ValidateLimit(req.GetLimit()); err != nil {
			violations = append(violations, FieldViolation("limit", err))
		}
	}

	if req.Offset != nil {
		if err := validator.ValidateOffset(req.GetOffset()); err != nil {
			violations = append(violations, FieldViolation("offset", err))
		}
	}

	return violations
}
