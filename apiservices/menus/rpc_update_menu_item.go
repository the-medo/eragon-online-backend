package menus

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/api/e"
	"github.com/the-medo/talebound-backend/apiservices/srv"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *ServiceMenus) UpdateMenuItem(ctx context.Context, req *pb.UpdateMenuItemRequest) (*pb.MenuItem, error) {
	violations := validateUpdateMenuItemRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	_, err := server.CheckMenuPermissions(ctx, req.GetMenuId(), &srv.ModulePermission{
		NeedsMenuPermission: true,
	})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed update menu: %v", err)
	}

	// if position is set, move the menu item to that position
	if req.Position != nil {
		positionChangeArg := db.MenuItemChangePositionsParams{
			MenuItemID:     req.GetMenuItemId(),
			TargetPosition: req.GetPosition(),
		}
		err = server.Store.MenuItemChangePositions(ctx, positionChangeArg)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to change menu item position: %s", err)
		}
	}

	// update the rest of the menu item
	arg := db.UpdateMenuItemParams{
		ID: req.GetMenuItemId(),
		MenuItemCode: sql.NullString{
			String: req.GetCode(),
			Valid:  req.Code != nil,
		},
		Name: sql.NullString{
			String: req.GetName(),
			Valid:  req.Name != nil,
		},
		IsMain: sql.NullBool{
			Bool:  req.GetIsMain(),
			Valid: req.IsMain != nil,
		},
		DescriptionPostID: sql.NullInt32{
			Int32: req.GetDescriptionPostId(),
			Valid: req.DescriptionPostId != nil,
		},
	}

	menuItem, err := server.Store.UpdateMenuItem(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update menu item: %s", err)
	}

	rsp := converters.ConvertMenuItem(menuItem)

	return rsp, nil
}

func validateUpdateMenuItemRequest(req *pb.UpdateMenuItemRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMenuId(req.GetMenuId()); err != nil {
		violations = append(violations, e.FieldViolation("menu_id", err))
	}

	if err := validator.ValidateMenuItemId(req.GetMenuItemId()); err != nil {
		violations = append(violations, e.FieldViolation("menu_item_id", err))
	}

	if req.Code != nil {
		if err := validator.ValidateMenuItemCode(req.GetCode()); err != nil {
			violations = append(violations, e.FieldViolation("code", err))
		}
	}

	if req.Name != nil {
		if err := validator.ValidateMenuItemName(req.GetName()); err != nil {
			violations = append(violations, e.FieldViolation("name", err))
		}
	}

	if req.Position != nil {
		if err := validator.ValidateMenuItemPosition(req.GetPosition()); err != nil {
			violations = append(violations, e.FieldViolation("position", err))
		}
	}

	if req.DescriptionPostId != nil {
		if err := validator.ValidatePostId(req.GetDescriptionPostId()); err != nil {
			violations = append(violations, e.FieldViolation("description_post_id", err))
		}
	}

	return violations
}
