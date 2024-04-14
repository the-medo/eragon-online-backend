package maps

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

func (server *ServiceMaps) CreateMapPinTypeGroup(ctx context.Context, request *pb.CreateMapPinTypeGroupRequest) (*pb.MapPinTypeGroup, error) {
	violations := validateCreateMapPinTypeGroup(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	_, err := server.CheckModuleIdPermissions(ctx, request.GetModuleId(), &servicecore.ModulePermission{
		NeedsEntityPermission: &[]db.EntityType{db.EntityTypeMap},
	})
	if err != nil {
		return nil, err
	}

	newPinTypeGroup, err := server.Store.CreateMapPinTypeGroup(ctx, request.GetName())
	if err != nil {
		return nil, err
	}

	argPinTypeGroup := db.CreateModuleMapPinTypeGroupParams{
		MapPinTypeGroupID: newPinTypeGroup.ID,
		ModuleID:          request.GetModuleId(),
	}

	_, err = server.Store.CreateModuleMapPinTypeGroup(ctx, argPinTypeGroup)
	if err != nil {
		return nil, err
	}

	rsp := converters.ConvertMapPinTypeGroup(newPinTypeGroup)

	return rsp, nil
}

func validateCreateMapPinTypeGroup(req *pb.CreateMapPinTypeGroupRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateUniversalName(req.GetName()); err != nil {
		violations = append(violations, e.FieldViolation("name", err))
	}

	return violations
}
