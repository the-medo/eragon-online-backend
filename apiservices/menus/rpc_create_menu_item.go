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

func (server *ServiceMenus) CreateMenuItem(ctx context.Context, req *pb.CreateMenuItemRequest) (*pb.MenuItem, error) {
	violations := validateCreateMenuItemRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	_, err := server.CheckMenuPermissions(ctx, req.GetMenuId(), &srv.ModulePermission{
		NeedsMenuPermission: true,
	})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed create menu item: %v", err)
	}

	argCreateEntityGroup := db.CreateEntityGroupParams{
		Name: sql.NullString{
			String: req.GetName(),
			Valid:  true,
		},
		Description: sql.NullString{},
		Style: sql.NullString{
			String: "framed",
			Valid:  true,
		},
		Direction: sql.NullString{
			String: "vertical",
			Valid:  true,
		},
	}

	newEntityGroup, err := server.Store.CreateEntityGroup(ctx, argCreateEntityGroup)
	if err != nil {
		return nil, err
	}

	arg := db.CreateMenuItemParams{
		MenuID:       req.GetMenuId(),
		MenuItemCode: req.GetCode(),
		Name:         req.GetName(),
		Position:     req.GetPosition(),
		IsMain: sql.NullBool{
			Bool:  req.GetIsMain(),
			Valid: req.IsMain != nil,
		},
		EntityGroupID: sql.NullInt32{
			Int32: newEntityGroup.ID,
			Valid: true,
		},
	}

	menuItem, err := server.Store.CreateMenuItem(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create menu item: %s", err)
	}

	rsp := converters.ConvertMenuItem(menuItem)

	return rsp, nil
}

func validateCreateMenuItemRequest(req *pb.CreateMenuItemRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMenuId(req.GetMenuId()); err != nil {
		violations = append(violations, e.FieldViolation("menu_id", err))
	}

	if err := validator.ValidateMenuItemCode(req.GetCode()); err != nil {
		violations = append(violations, e.FieldViolation("code", err))
	}

	if err := validator.ValidateMenuItemName(req.GetName()); err != nil {
		violations = append(violations, e.FieldViolation("name", err))
	}

	if req.Position != nil {
		if err := validator.ValidateMenuItemPosition(req.GetPosition()); err != nil {
			violations = append(violations, e.FieldViolation("position", err))
		}
	}

	return violations
}
