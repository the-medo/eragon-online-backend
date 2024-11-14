package systems

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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *ServiceSystems) UpdateSystem(ctx context.Context, req *pb.UpdateSystemRequest) (*pb.System, error) {
	violations := validateUpdateSystemRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	var needsEntityPermission []db.EntityType

	_, _, err := server.CheckModuleTypePermissions(ctx, db.ModuleTypeSystem, req.GetSystemId(), &servicecore.ModulePermission{
		NeedsSuperAdmin:       true,
		NeedsEntityPermission: &needsEntityPermission,
	})

	if err != nil {
		return nil, err
	}

	arg := db.UpdateSystemParams{
		SystemID:         req.GetSystemId(),
		Name:             sql.NullString{String: req.GetName(), Valid: req.Name != nil},
		ShortDescription: sql.NullString{String: req.GetShortDescription(), Valid: req.ShortDescription != nil},
	}

	system, err := server.Store.UpdateSystem(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update system: %s", err)
	}

	return converters.ConvertSystem(system), nil
}

func validateUpdateSystemRequest(req *pb.UpdateSystemRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateModuleId(req.GetSystemId()); err != nil {
		violations = append(violations, e.FieldViolation("system_id", err))
	}

	if req.Name != nil {
		if err := validator.ValidateModuleName(req.GetName()); err != nil {
			violations = append(violations, e.FieldViolation("name", err))
		}
	}

	if req.ShortDescription != nil {
		if err := validator.ValidateModuleShortDescription(req.GetShortDescription()); err != nil {
			violations = append(violations, e.FieldViolation("short_description", err))
		}
	}

	return violations
}
