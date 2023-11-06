package modules

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/api/e"
	"github.com/the-medo/talebound-backend/apiservices/srv"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *ServiceModules) DeleteModuleAdmin(ctx context.Context, req *pb.DeleteModuleAdminRequest) (*emptypb.Empty, error) {
	violations := validateDeleteModuleAdmin(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	authPayload, err := server.AuthorizeUserCookie(ctx)
	if err != nil {
		return nil, e.UnauthenticatedError(err)
	}
	//user can remove himself from module even if he is not super admin
	_, err = server.CheckModuleIdPermissions(ctx, req.GetModuleId(), &srv.ModulePermission{
		NeedsSuperAdmin: req.GetUserId() != authPayload.UserId,
	})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to delete module admin: %v", err)
	}

	arg := db.DeleteModuleAdminParams{
		ModuleID: req.GetModuleId(),
		UserID:   req.GetUserId(),
	}

	err = server.Store.DeleteModuleAdmin(ctx, arg)
	if err != nil {
		return nil, err
	}

	userModule, err := server.Store.UpsertUserModule(ctx, db.UpsertUserModuleParams{
		ModuleID: req.GetModuleId(),
		UserID:   req.GetUserId(),
		Admin: sql.NullBool{
			Bool:  false,
			Valid: true,
		},
	})

	if userModule.Admin == false && userModule.Following == false && userModule.Favorite == false {
		err = server.Store.DeleteUserModule(ctx, db.DeleteUserModuleParams{
			ModuleID: req.GetModuleId(),
			UserID:   req.GetUserId(),
		})
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func validateDeleteModuleAdmin(req *pb.DeleteModuleAdminRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateUserId(req.GetUserId()); err != nil {
		violations = append(violations, e.FieldViolation("user_id", err))
	}

	if err := validator.ValidateUniversalId(req.GetModuleId()); err != nil {
		violations = append(violations, e.FieldViolation("module_id", err))
	}

	return violations
}
