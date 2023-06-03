package api

import (
	"context"
	"database/sql"
	"fmt"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.Post, error) {
	violations := validateCreatePostRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUserCookie(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	arg := db.CreatePostParams{
		UserID:     authPayload.UserId,
		PostTypeID: req.GetPostTypeId(),
		Title:      req.GetTitle(),
		Content:    req.GetContent(),
	}

	postResult, err := server.store.CreatePost(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create post: %s", err)
	}

	rsp := convertPost(postResult)

	return rsp, nil
}

func (server *Server) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.Post, error) {
	violations := validateUpdatePostRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	authPayload, err := server.authorizeUserCookie(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	_, err = server.store.InsertPostHistory(ctx, req.GetPostId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert post history: %v", err)
	}

	arg := db.UpdatePostParams{
		PostID: req.GetPostId(),
		Title: sql.NullString{
			String: req.GetTitle(),
			Valid:  req.Title != nil,
		},
		Content: sql.NullString{
			String: req.GetContent(),
			Valid:  req.Title != nil,
		},
		PostTypeID: sql.NullInt32{
			Int32: req.GetPostTypeId(),
			Valid: req.PostTypeId != nil,
		},
		LastUpdatedUserID: sql.NullInt32{
			Int32: authPayload.UserId,
			Valid: true,
		},
	}

	post, err := server.store.UpdatePost(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update post: %v", err)
	}
	rsp := convertPost(post)

	return rsp, nil
}

func (server *Server) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	violations := validateDeletePostRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUserCookie(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	post, err := server.store.GetPostById(ctx, req.GetPostId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get post: %v", err)
	}

	if post.UserID != authPayload.UserId {
		err := server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
		if err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "can not update this post - you are not creator or admin: %v", err)
		}
	}

	err = server.store.DeletePost(ctx, req.GetPostId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete post: %v", err)
	}

	return &pb.DeletePostResponse{
		Success: true,
		Message: fmt.Sprintf("Post %s (%v) deleted successfully.", post.Title, req.GetPostId()),
	}, nil
}

// ================= VALIDATION =================

func validateCreatePostRequest(req *pb.CreatePostRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidatePostTitle(req.GetTitle()); err != nil {
		violations = append(violations, FieldViolation("title", err))
	}

	return violations
}

func validateUpdatePostRequest(req *pb.UpdatePostRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidatePostId(req.GetPostId()); err != nil {
		violations = append(violations, FieldViolation("post_id", err))
	}

	if req.Title != nil {
		if err := validator.ValidatePostTitle(req.GetTitle()); err != nil {
			violations = append(violations, FieldViolation("title", err))
		}
	}

	if req.PostTypeId != nil {
		if err := validator.ValidatePostTypeId(req.GetPostTypeId()); err != nil {
			violations = append(violations, FieldViolation("post_type_id", err))
		}
	}

	return violations
}

func validateDeletePostRequest(req *pb.DeletePostRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidatePostId(req.GetPostId()); err != nil {
		violations = append(violations, FieldViolation("post_id", err))
	}

	return violations
}
