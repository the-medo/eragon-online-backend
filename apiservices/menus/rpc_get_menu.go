package menus

import (
	"context"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/api/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceMenus) GetMenu(ctx context.Context, req *pb.GetMenuRequest) (*pb.ViewMenu, error) {
	violations := validateGetMenuRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	menu, err := server.Store.GetMenu(ctx, req.MenuId)
	if err != nil {
		return nil, err
	}

	rsp := converters.ConvertViewMenu(menu)

	return rsp, nil
}

func validateGetMenuRequest(req *pb.GetMenuRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMenuId(req.GetMenuId()); err != nil {
		violations = append(violations, e.FieldViolation("menu_id", err))
	}

	return violations
}
