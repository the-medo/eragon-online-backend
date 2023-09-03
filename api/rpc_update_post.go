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
			Valid:  req.Content != nil,
		},
		Description: sql.NullString{
			String: req.GetDescription(),
			Valid:  req.Description != nil,
		},
		PostTypeID: sql.NullInt32{
			Int32: req.GetPostTypeId(),
			Valid: req.PostTypeId != nil,
		},
		LastUpdatedUserID: sql.NullInt32{
			Int32: authPayload.UserId,
			Valid: true,
		},
		IsDraft: sql.NullBool{
			Bool:  req.GetIsDraft(),
			Valid: req.IsDraft != nil,
		},
		IsPrivate: sql.NullBool{
			Bool:  req.GetIsPrivate(),
			Valid: req.IsPrivate != nil,
		},
		ThumbnailImgID: sql.NullInt32{
			Int32: req.GetImageThumbnailId(),
			Valid: req.ImageThumbnailId != nil,
		},
	}

	post, err := server.store.UpdatePost(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update post: %v", err)
	}

	postType, err := server.store.GetPostTypeById(ctx, post.PostTypeID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get post type: %v", err)
	}

	rsp := convertPostAndPostType(post, postType)

	return rsp, nil
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

	if req.Description != nil {
		if err := validator.ValidatePostDescription(req.GetDescription()); err != nil {
			violations = append(violations, FieldViolation("description", err))
		}
	}

	if req.Content != nil {
		if err := validator.ValidatePostContent(req.GetContent()); err != nil {
			violations = append(violations, FieldViolation("content", err))
		}
	}

	if req.PostTypeId != nil {
		if err := validator.ValidatePostTypeId(req.GetPostTypeId()); err != nil {
			violations = append(violations, FieldViolation("post_type_id", err))
		}
	}

	if req.ImageThumbnailId != nil {
		if err := validator.ValidateImageId(req.GetImageThumbnailId()); err != nil {
			violations = append(violations, FieldViolation("image_thumbnail_id", err))
		}
	}

	return violations
}
