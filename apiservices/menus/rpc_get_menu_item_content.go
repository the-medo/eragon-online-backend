package menus

import (
	"context"
	"github.com/the-medo/talebound-backend/api/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceMenus) GetMenuItemContent(ctx context.Context, req *pb.GetMenuItemContentRequest) (*pb.GetMenuItemContentResponse, error) {
	violations := validateGetMenuItemContentRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	rsp := &pb.GetMenuItemContentResponse{}

	return rsp, nil
}

func validateGetMenuItemContentRequest(req *pb.GetMenuItemContentRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMenuId(req.GetMenuId()); err != nil {
		violations = append(violations, e.FieldViolation("menu_id", err))
	}
	if err := validator.ValidateMenuItemId(req.GetMenuItemId()); err != nil {
		violations = append(violations, e.FieldViolation("menu_item_id", err))
	}

	return violations
}
