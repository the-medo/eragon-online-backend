package menus

import (
	"context"
	"github.com/the-medo/talebound-backend/api/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *ServiceMenus) UpdateMenuItemMoveGroupUp(ctx context.Context, req *pb.UpdateMenuItemMoveGroupUpRequest) (*emptypb.Empty, error) {
	violations := validateUpdateMenuItemMoveGroupUpRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	_, err := server.CheckMenuAdmin(ctx, req.GetMenuId(), false)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to move menu item group: %v", err)
	}

	err = server.Store.MenuItemMoveGroupUp(ctx, req.MenuItemId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to move menu item group: %s", err)
	}

	return nil, nil
}

func validateUpdateMenuItemMoveGroupUpRequest(req *pb.UpdateMenuItemMoveGroupUpRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMenuId(req.GetMenuId()); err != nil {
		violations = append(violations, e.FieldViolation("menu_id", err))
	}

	if err := validator.ValidateMenuItemId(req.GetMenuItemId()); err != nil {
		violations = append(violations, e.FieldViolation("menu_item_id", err))
	}
	return violations
}
