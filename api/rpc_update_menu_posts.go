package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *Server) UpdateMenuPosts(ctx context.Context, req *pb.UpdateMenuPostsRequest) (*emptypb.Empty, error) {
	violations := validateUpdateMenuPostsRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	_, err := server.CheckMenuAdmin(ctx, req.GetMenuId(), false)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed update menu item post: %v", err)
	}

	postIDs := req.GetPostIds()
	for _, id := range postIDs {
		// delete post id

		argDelete := db.DeleteMenuItemPostParams{
			MenuID: req.GetMenuId(),
			PostID: id,
		}

		err = server.store.DeleteMenuItemPost(ctx, argDelete)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to delete menu item post: %s", err)
		}

		// create new post id

		argCreate := db.CreateMenuItemPostParams{
			MenuID: req.GetMenuId(),
			PostID: id,
		}

		if req.GetMenuItemId() != 0 {
			argCreate.MenuItemID = sql.NullInt32{
				Int32: req.GetMenuItemId(),
				Valid: true,
			}
		}

		_, err := server.store.CreateMenuItemPost(ctx, argCreate)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create menu item post [CreateMenuItemPost]: %s", err)
		}
	}

	return nil, nil
}

func validateUpdateMenuPostsRequest(req *pb.UpdateMenuPostsRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateMenuId(req.GetMenuId()); err != nil {
		violations = append(violations, FieldViolation("menu_id", err))
	}

	if req.MenuItemId != nil {
		if err := validator.ValidateMenuItemId(req.GetMenuItemId()); err != nil {
			violations = append(violations, FieldViolation("menu_item_id", err))
		}
	}

	postIDs := req.GetPostIds()
	if len(postIDs) == 0 {
		violations = append(violations, FieldViolation("post_ids", errors.New("post IDs cannot be empty")))
	}

	// Use a map to track seen IDs for detecting duplicates.
	seen := make(map[int32]bool)
	for _, id := range postIDs {

		if err := validator.ValidatePostId(id); err != nil {
			violations = append(violations, FieldViolation("post_ids", err))
		}

		if seen[id] {
			violations = append(violations, FieldViolation("post_ids", fmt.Errorf("duplicate post ID: %d", id)))
		}
		seen[id] = true
	}

	return violations
}
