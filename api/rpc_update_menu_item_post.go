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
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *Server) UpdateMenuItemPost(ctx context.Context, req *pb.UpdateMenuItemPostRequest) (*emptypb.Empty, error) {
	violations := validateUpdateMenuItemPostRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	_, err := server.CheckMenuAdmin(ctx, req.GetMenuId(), false)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed update menu item post: %v", err)
	}

	// if position is set, move the menu item to that position
	if req.Position != nil {
		positionChangeArg := db.MenuItemPostChangePositionsParams{
			MenuItemID:     req.GetMenuItemId(),
			PostID:         req.GetPostId(),
			TargetPosition: req.GetPosition(),
		}
		err = server.store.MenuItemPostChangePositions(ctx, positionChangeArg)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to change menu item post position: %s", err)
		}
	}

	if req.NewMenuItemId != nil {
		if req.GetNewMenuItemId() > 0 {
			updateArg := db.UpdateMenuItemPostParams{
				MenuItemID: sql.NullInt32{
					Int32: req.GetMenuItemId(),
					Valid: true,
				},
				PostID: sql.NullInt32{
					Int32: req.GetPostId(),
					Valid: true,
				},
				NewMenuItemID: sql.NullInt32{
					Int32: req.GetNewMenuItemId(),
					Valid: true,
				},
			}

			_, err = server.store.UpdateMenuItemPost(ctx, updateArg)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to update menu item post: %s", err)
			}
		} else {

			updateArg := db.UnassignMenuItemPostParams{
				MenuItemID: sql.NullInt32{
					Int32: req.GetMenuItemId(),
					Valid: true,
				},
				PostID: req.GetPostId(),
			}
			_, err = server.store.UnassignMenuItemPost(ctx, updateArg)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to unassign menu item post: %s", err)
			}
		}
	}

	return nil, nil
}

func validateUpdateMenuItemPostRequest(req *pb.UpdateMenuItemPostRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMenuId(req.GetMenuId()); err != nil {
		violations = append(violations, FieldViolation("menu_id", err))
	}

	if err := validator.ValidateMenuItemId(req.GetMenuItemId()); err != nil {
		violations = append(violations, FieldViolation("menu_item_id", err))
	}

	if err := validator.ValidatePostId(req.GetPostId()); err != nil {
		violations = append(violations, FieldViolation("post_id", err))
	}

	if req.Position != nil {
		if err := validator.ValidateMenuItemPosition(req.GetPosition()); err != nil {
			violations = append(violations, FieldViolation("position", err))
		}
	}

	return violations
}
