package maps

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

func (server *ServiceMaps) UpdateMapPinTypeGroup(ctx context.Context, request *pb.UpdateMapPinTypeGroupRequest) (*pb.MapPinTypeGroup, error) {
	violations := validateUpdateMapPinTypeGroup(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	_, err := server.CheckModuleIdPermissions(ctx, request.GetModuleId(), &servicecore.ModulePermission{
		NeedsEntityPermission: &[]db.EntityType{db.EntityTypeMap},
	})

	if err != nil {
		return nil, err
	}

	argMapPinTypeGroup := db.UpdateMapPinTypeGroupParams{
		ID:   request.GetMapPinTypeGroupId(),
		Name: sql.NullString{String: request.GetName(), Valid: true},
	}

	updatedMapPinTypeGroup, err := server.Store.UpdateMapPinTypeGroup(ctx, argMapPinTypeGroup)
	if err != nil {
		return nil, err
	}

	rsp := converters.ConvertMapPinTypeGroup(updatedMapPinTypeGroup)

	return rsp, nil
}

func validateUpdateMapPinTypeGroup(req *pb.UpdateMapPinTypeGroupRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUniversalId(req.GetModuleId()); err != nil {
		violations = append(violations, e.FieldViolation("module_id", err))
	}

	if err := validator.ValidateUniversalId(req.GetMapPinTypeGroupId()); err != nil {
		violations = append(violations, e.FieldViolation("pin_type_id", err))
	}

	if err := validator.ValidateUniversalName(req.GetName()); err != nil {
		violations = append(violations, e.FieldViolation("name", err))
	}

	return violations
}
