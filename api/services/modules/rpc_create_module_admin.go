package modules

import (
	"context"
	"github.com/the-medo/talebound-backend/api/servicecore"
	"github.com/the-medo/talebound-backend/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceModules) CreateModuleAdmin(ctx context.Context, request *pb.CreateModuleAdminRequest) (*pb.ModuleAdmin, error) {
	violations := validateCreateModuleAdmin(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	_, err := server.CheckModuleIdPermissions(ctx, request.GetModuleId(), &servicecore.ModulePermission{
		NeedsSuperAdmin: true,
	})
	if err != nil {
		return nil, err
	}

	authPayload, err := server.AuthorizeUserCookie(ctx)
	if err != nil {
		return nil, e.UnauthenticatedError(err)
	}

	arg := db.CreateModuleAdminParams{
		ModuleID:           request.GetModuleId(),
		UserID:             authPayload.UserId,
		SuperAdmin:         false,
		Approved:           2,
		MotivationalLetter: request.GetMotivationalLetter(),
	}

	moduleAdmin, err := server.Store.CreateModuleAdmin(ctx, arg)
	if err != nil {
		return nil, err
	}

	user, err := server.Store.GetUserById(ctx, authPayload.UserId)

	rsp := converters.ConvertModuleAdmin(moduleAdmin, user)

	return rsp, nil
}

func validateCreateModuleAdmin(req *pb.CreateModuleAdminRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUniversalId(req.GetModuleId()); err != nil {
		violations = append(violations, e.FieldViolation("module_id", err))
	}

	if err := validator.ValidateModuleAdminMotivationalLetter(req.GetMotivationalLetter()); err != nil {
		violations = append(violations, e.FieldViolation("motivational_letter", err))
	}

	return violations
}
