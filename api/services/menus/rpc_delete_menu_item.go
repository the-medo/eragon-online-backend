package menus

import (
	"context"
	"github.com/the-medo/talebound-backend/api/servicecore"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *ServiceMenus) DeleteMenuItem(ctx context.Context, req *pb.DeleteMenuItemRequest) (*emptypb.Empty, error) {
	violations := validateDeleteMenuItemRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	_, err := server.CheckMenuPermissions(ctx, req.GetMenuId(), &servicecore.ModulePermission{
		NeedsMenuPermission: true,
	})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed delete menu item: %v", err)
	}

	err = server.Store.DeleteMenuItem(ctx, req.GetMenuItemId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete menu item: %s", err)
	}

	return nil, nil
}

func validateDeleteMenuItemRequest(req *pb.DeleteMenuItemRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMenuId(req.GetMenuId()); err != nil {
		violations = append(violations, e.FieldViolation("menu_id", err))
	}

	if err := validator.ValidateMenuItemId(req.GetMenuItemId()); err != nil {
		violations = append(violations, e.FieldViolation("menu_item_id", err))
	}

	return violations
}
