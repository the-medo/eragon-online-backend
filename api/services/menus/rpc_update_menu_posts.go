package menus

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/the-medo/talebound-backend/api/servicecore"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *ServiceMenus) UpdateMenuPosts(ctx context.Context, req *pb.UpdateMenuPostsRequest) (*pb.UpdateMenuPostsResponse, error) {
	violations := validateUpdateMenuPostsRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	_, err := server.CheckMenuPermissions(ctx, req.GetMenuId(), &servicecore.ModulePermission{
		NeedsMenuPermission: true,
	})
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

		err = server.Store.DeleteMenuItemPost(ctx, argDelete)
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

		_, err := server.Store.CreateMenuItemPost(ctx, argCreate)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create menu item post [CreateMenuItemPost]: %s", err)
		}
	}

	rsp, err := server.GetMenuItemPostsByMenuId(ctx, &pb.GetMenuItemPostsByMenuIdRequest{
		MenuId: req.GetMenuId(),
	})

	return &pb.UpdateMenuPostsResponse{MenuItemPosts: rsp.MenuItemPosts}, nil
}

func validateUpdateMenuPostsRequest(req *pb.UpdateMenuPostsRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateMenuId(req.GetMenuId()); err != nil {
		violations = append(violations, e.FieldViolation("menu_id", err))
	}

	if req.MenuItemId != nil {
		if err := validator.ValidateMenuItemId(req.GetMenuItemId()); err != nil {
			violations = append(violations, e.FieldViolation("menu_item_id", err))
		}
	}

	postIDs := req.GetPostIds()
	if len(postIDs) == 0 {
		violations = append(violations, e.FieldViolation("post_ids", errors.New("post IDs cannot be empty")))
	}

	// Use a map to track seen IDs for detecting duplicates.
	seen := make(map[int32]bool)
	for _, id := range postIDs {

		if err := validator.ValidatePostId(id); err != nil {
			violations = append(violations, e.FieldViolation("post_ids", err))
		}

		if seen[id] {
			violations = append(violations, e.FieldViolation("post_ids", fmt.Errorf("duplicate post ID: %d", id)))
		}
		seen[id] = true
	}

	return violations
}
