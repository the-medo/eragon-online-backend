package menus

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/api/e"
	"github.com/the-medo/talebound-backend/apiservices/servicecore"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *ServiceMenus) UpdateMenu(ctx context.Context, req *pb.UpdateMenuRequest) (*pb.ViewMenu, error) {
	violations := validateUpdateMenu(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	_, err := server.CheckMenuPermissions(ctx, req.GetMenuId(), &servicecore.ModulePermission{
		NeedsMenuPermission: true,
	})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed update menu: %v", err)
	}

	arg := db.UpdateMenuParams{
		ID: req.GetMenuId(),
		MenuCode: sql.NullString{
			String: req.GetCode(),
			Valid:  req.Code != nil,
		},
		MenuHeaderImgID: sql.NullInt32{
			Int32: req.GetHeaderImgId(),
			Valid: req.HeaderImgId != nil,
		},
	}

	_, err = server.Store.UpdateMenu(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to update menu: %s", err)
	}

	menu, err := server.Store.GetMenu(ctx, req.GetMenuId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get menu: %s", err)
	}

	rsp := converters.ConvertViewMenu(menu)

	return rsp, nil
}

func validateUpdateMenu(req *pb.UpdateMenuRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMenuId(req.GetMenuId()); err != nil {
		violations = append(violations, e.FieldViolation("menu_id", err))
	}

	if req.Code != nil {
		if err := validator.ValidateMenuCode(req.GetCode()); err != nil {
			violations = append(violations, e.FieldViolation("code", err))
		}
	}

	if req.HeaderImgId != nil {
		if err := validator.ValidateImageId(req.GetHeaderImgId()); err != nil {
			violations = append(violations, e.FieldViolation("header_img_id", err))
		}
	}

	return violations
}
