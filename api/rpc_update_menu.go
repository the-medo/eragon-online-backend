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

func (server *Server) UpdateMenu(ctx context.Context, req *pb.UpdateMenuRequest) (*pb.Menu, error) {
	violations := validateUpdateMenuRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	_, err := server.CheckMenuAdmin(ctx, req.GetMenuId(), false)
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

	menu, err := server.store.UpdateMenu(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to update menu: %s", err)
	}

	rsp := converters.ConvertMenu(menu)

	return rsp, nil
}

func validateUpdateMenuRequest(req *pb.UpdateMenuRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateInt(req.GetMenuId(), 1, 4000); err != nil {
		violations = append(violations, FieldViolation("menu_id", err))
	}

	return violations
}
