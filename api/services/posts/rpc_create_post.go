package posts

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/api/servicecore"
	"github.com/the-medo/talebound-backend/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *ServicePosts) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.Post, error) {
	violations := validateCreatePostRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	_, err := server.CheckModuleIdPermissions(ctx, req.GetModuleId(), &servicecore.ModulePermission{
		NeedsEntityPermission: &[]db.EntityType{db.EntityTypePost},
	})

	authPayload, err := server.AuthorizeUserCookie(ctx)
	if err != nil {
		return nil, e.UnauthenticatedError(err)
	}

	arg := db.CreatePostParams{
		UserID: authPayload.UserId,
		Title:  req.GetTitle(),
		Description: sql.NullString{
			String: req.GetDescription(),
			Valid:  req.Description != nil,
		},
		Content:   req.GetContent(),
		IsDraft:   req.GetIsDraft(),
		IsPrivate: req.GetIsPrivate(),
		ThumbnailImgID: sql.NullInt32{
			Int32: req.GetImageThumbnailId(),
			Valid: req.ImageThumbnailId != nil,
		},
	}

	postResult, err := server.Store.CreatePost(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create post: %s", err)
	}

	_, err = server.Store.CreateEntity(ctx, db.CreateEntityParams{
		Type:     db.EntityTypePost,
		ModuleID: req.GetModuleId(),
		PostID: sql.NullInt32{
			Int32: postResult.ID,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	post, err := server.Store.GetPostById(ctx, postResult.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get post: %s", err)
	}

	rsp := converters.ConvertPost(post)

	return rsp, nil
}

func validateCreatePostRequest(req *pb.CreatePostRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUniversalId(req.GetModuleId()); err != nil {
		violations = append(violations, e.FieldViolation("module_id", err))
	}

	if err := validator.ValidatePostTitle(req.GetTitle()); err != nil {
		violations = append(violations, e.FieldViolation("title", err))
	}

	if err := validator.ValidatePostContent(req.GetContent()); err != nil {
		violations = append(violations, e.FieldViolation("content", err))
	}

	return violations
}
