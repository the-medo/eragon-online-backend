package modules

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/api/servicecore"
	"github.com/the-medo/talebound-backend/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceModules) UpdateModuleAdmin(ctx context.Context, req *pb.UpdateModuleAdminRequest) (*pb.ModuleAdmin, error) {
	violations := validateUpdateModuleAdmin(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	_, err := server.CheckModuleIdPermissions(ctx, req.GetModuleId(), &servicecore.ModulePermission{
		NeedsSuperAdmin: true,
	})
	if err != nil {
		return nil, err
	}

	//TODO: add AllowedEntityTypes and AllowedMenu here and to SQL query
	arg := db.UpdateModuleAdminParams{
		ModuleID: req.GetModuleId(),
		UserID:   req.GetUserId(),

		SuperAdmin: sql.NullBool{
			Bool:  req.GetSuperAdmin(),
			Valid: req.SuperAdmin != nil,
		},
		Approved: sql.NullInt32{
			Int32: req.GetApproved(),
			Valid: req.Approved != nil,
		},
		MotivationalLetter: sql.NullString{
			String: req.GetMotivationalLetter(),
			Valid:  req.MotivationalLetter != nil,
		},
		AllowedEntityTypesPresent: req.GetAllowedEntityTypes() != nil,
		AllowedEntityTypes:        make([]db.EntityType, 0),
		AllowedMenu: sql.NullBool{
			Bool:  req.GetSuperAdmin(),
			Valid: req.SuperAdmin != nil,
		},
	}

	if req.GetAllowedEntityTypes() != nil {
		for _, et := range req.AllowedEntityTypes.EntityTypes {
			arg.AllowedEntityTypes = append(arg.AllowedEntityTypes, converters.ConvertEntityTypeToDB(et))
		}
	}

	moduleAdmin, err := server.Store.UpdateModuleAdmin(ctx, arg)
	if err != nil {
		return nil, err
	}

	user, err := server.Store.GetUserById(ctx, req.GetUserId())

	rsp := converters.ConvertModuleAdmin(moduleAdmin, user)

	return rsp, nil
}

func validateUpdateModuleAdmin(req *pb.UpdateModuleAdminRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateUserId(req.GetUserId()); err != nil {
		violations = append(violations, e.FieldViolation("user_id", err))
	}

	if err := validator.ValidateUniversalId(req.GetModuleId()); err != nil {
		violations = append(violations, e.FieldViolation("module_id", err))
	}

	if req.Approved != nil {
		if err := validator.ValidateModuleAdminApproved(req.GetApproved()); err != nil {
			violations = append(violations, e.FieldViolation("approved", err))
		}
	}

	if req.MotivationalLetter != nil {
		if err := validator.ValidateModuleAdminMotivationalLetter(req.GetMotivationalLetter()); err != nil {
			violations = append(violations, e.FieldViolation("motivational_letter", err))
		}
	}

	return violations
}
