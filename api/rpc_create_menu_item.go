package api

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/api/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateMenuItem(ctx context.Context, req *pb.CreateMenuItemRequest) (*pb.MenuItem, error) {
	violations := validateCreateMenuItemRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	_, err := server.CheckMenuAdmin(ctx, req.GetMenuId(), false)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed create menu item: %v", err)
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
	}

	menuItem, err := server.store.CreateMenuItem(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create menu item: %s", err)
	}

	rsp := converters.ConvertMenuItem(menuItem)

	return rsp, nil
}

func validateCreateMenuItemRequest(req *pb.CreateMenuItemRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMenuId(req.GetMenuId()); err != nil {
		violations = append(violations, FieldViolation("menu_id", err))
	}

	if err := validator.ValidateMenuItemCode(req.GetCode()); err != nil {
		violations = append(violations, FieldViolation("code", err))
	}

	if err := validator.ValidateMenuItemName(req.GetName()); err != nil {
		violations = append(violations, FieldViolation("name", err))
	}

	if req.Position != nil {
		if err := validator.ValidateMenuItemPosition(req.GetPosition()); err != nil {
			violations = append(violations, FieldViolation("position", err))
		}
	}

	return violations
}
