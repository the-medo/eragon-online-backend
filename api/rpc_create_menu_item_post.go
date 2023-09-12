package api

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/consts"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateMenuItemPost(ctx context.Context, req *pb.CreateMenuItemPostRequest) (*pb.MenuItemPost, error) {
	violations := validateCreateMenuItemPostRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	auth, err := server.CheckMenuAdmin(ctx, req.GetMenuId(), false)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed create menu item post [permission]: %v", err)
	}

	isDescriptionPost := req.GetIsMenuItemDescriptionPost()

	argPost := db.CreatePostParams{
		UserID:     auth.UserId,
		Title:      "Title",
		PostTypeID: consts.PostTypeWorldDescription,
		IsDraft:    true,
		IsPrivate:  false,
	}

	if req.Title != nil && req.GetTitle() != "" {
		argPost.Title = req.GetTitle()
	}

	if req.ShortDescription != nil && req.GetShortDescription() != "" {
		argPost.Description = sql.NullString{
			String: req.GetShortDescription(),
			Valid:  true,
		}
	}

	if req.ImageThumbnailId != nil {
		argPost.ThumbnailImgID = sql.NullInt32{
			Int32: req.GetImageThumbnailId(),
			Valid: true,
		}
	}

	newPost, err := server.store.CreatePost(ctx, argPost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create menu item post [CreatePost]: %s", err)
	}

	post, err := server.store.GetPostById(ctx, newPost.ID)

	rsp := &pb.MenuItemPost{
		MenuItemId: req.GetMenuItemId(),
		PostId:     post.ID,
		Position:   req.GetPosition(),
		Post:       converters.ConvertViewPost(post),
	}

	if isDescriptionPost {
		_, err := server.store.UpdateMenuItem(ctx, db.UpdateMenuItemParams{
			ID: req.GetMenuItemId(),
			DescriptionPostID: sql.NullInt32{
				Int32: newPost.ID,
				Valid: true,
			},
		})

		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update menu item post [UpdateMenuItem]: %s", err)
		}

	} else {
		arg := db.CreateMenuItemPostParams{
			MenuID: req.GetMenuId(),
			MenuItemID: sql.NullInt32{
				Int32: req.GetMenuItemId(),
				Valid: true,
			},
			PostID:   newPost.ID,
			Position: req.GetPosition(),
		}

		menuItemPost, err := server.store.CreateMenuItemPost(ctx, arg)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create menu item post [CreateMenuItemPost]: %s", err)
		}

		rsp = converters.ConvertMenuItemPost(menuItemPost, post)
		return rsp, nil
	}

	return rsp, nil
}

func validateCreateMenuItemPostRequest(req *pb.CreateMenuItemPostRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMenuId(req.GetMenuId()); err != nil {
		violations = append(violations, FieldViolation("menu_id", err))
	}

	if err := validator.ValidateMenuItemId(req.GetMenuItemId()); err != nil {
		violations = append(violations, FieldViolation("menu_item_id", err))
	}

	if req.PostId != nil {
		if err := validator.ValidatePostId(req.GetPostId()); err != nil {
			violations = append(violations, FieldViolation("post_id", err))
		}
	}

	if req.Position != nil {
		if err := validator.ValidateMenuItemPosition(req.GetPosition()); err != nil {
			violations = append(violations, FieldViolation("position", err))
		}
	}

	if req.Title != nil {
		if err := validator.ValidatePostTitle(req.GetTitle()); err != nil {
			violations = append(violations, FieldViolation("title", err))
		}
	}

	if req.ShortDescription != nil {
		if err := validator.ValidatePostDescription(req.GetShortDescription()); err != nil {
			violations = append(violations, FieldViolation("short_description", err))
		}
	}

	if req.ImageThumbnailId != nil {
		if err := validator.ValidateImageId(req.GetImageThumbnailId()); err != nil {
			violations = append(violations, FieldViolation("image_thumbnail_id", err))
		}
	}

	return violations
}
