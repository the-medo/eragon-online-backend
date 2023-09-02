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

func (server *Server) DeleteMenuItemPost(ctx context.Context, req *pb.DeleteMenuItemPostRequest) (*emptypb.Empty, error) {
	violations := validateDeleteMenuItemPostRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	_, err := server.CheckMenuAdmin(ctx, req.GetMenuId(), false)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed delete menu item post: %v", err)
	}

	arg := db.DeleteMenuItemPostParams{
		MenuItemID: sql.NullInt32{
			Int32: req.GetMenuItemId(),
			Valid: true,
		},
		PostID: req.GetPostId(),
	}

	err = server.store.DeleteMenuItemPost(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete menu item post: %s", err)
	}

	return nil, nil
}

func validateDeleteMenuItemPostRequest(req *pb.DeleteMenuItemPostRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMenuId(req.GetMenuId()); err != nil {
		violations = append(violations, FieldViolation("menu_id", err))
	}

	if err := validator.ValidateMenuItemId(req.GetMenuItemId()); err != nil {
		violations = append(violations, FieldViolation("menu_item_id", err))
	}

	if err := validator.ValidatePostId(req.GetPostId()); err != nil {
		violations = append(violations, FieldViolation("post_id", err))
	}

	return violations
}
