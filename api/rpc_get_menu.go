package api

import (
	"context"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *Server) GetMenu(ctx context.Context, req *pb.GetMenuRequest) (*pb.Menu, error) {
	violations := validateGetMenuRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	menu, err := server.store.GetMenu(ctx, req.MenuId)
	if err != nil {
		return nil, err
	}

	rsp := converters.ConvertMenu(menu)

	return rsp, nil
}

func validateGetMenuRequest(req *pb.GetMenuRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateInt(req.GetMenuId(), 1, 4000); err != nil {
		violations = append(violations, FieldViolation("menu_id", err))
	}

	return violations
}
