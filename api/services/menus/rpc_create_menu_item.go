package menus

import (
	"context"
	"github.com/the-medo/talebound-backend/api/servicecore"
	"github.com/the-medo/talebound-backend/converters"
	"github.com/the-medo/talebound-backend/e"
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

	_, err := server.CheckMenuPermissions(ctx, req.GetMenuId(), &servicecore.ModulePermission{
		NeedsMenuPermission: true,
	})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed create menu item: %v", err)
	}

	dbMenuItem, err := server.SharedCreateMenuItem(ctx, req)
	if err != nil {
		return nil, err
	}
	return converters.ConvertMenuItem(*dbMenuItem), nil
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
