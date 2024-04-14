package maps

import (
	"context"
	"fmt"
	"github.com/the-medo/talebound-backend/api/servicecore"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *ServiceMaps) DeleteMapPinTypeGroup(ctx context.Context, request *pb.DeleteMapPinTypeGroupRequest) (*emptypb.Empty, error) {
	violations := validateDeleteMapPinTypeGroup(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	_, err := server.CheckModuleIdPermissions(ctx, request.GetModuleId(), &servicecore.ModulePermission{
		NeedsEntityPermission: &[]db.EntityType{db.EntityTypeMap},
	})

	if err != nil {
		return nil, err
	}

	pinTypeRows, err := server.Store.GetMapPinTypesForModule(ctx, request.GetModuleId())
	if err != nil {
		return nil, err
	}

	for _, pinTypeRow := range pinTypeRows {
		if pinTypeRow.MapPinTypeGroupID == request.GetMapPinTypeGroupId() {
			return nil, fmt.Errorf("map_pin_type_group is not empty - can not delete")
		}
	}

	err = server.Store.DeleteModuleMapPinTypeGroup(ctx, db.DeleteModuleMapPinTypeGroupParams{
		ModuleID:          request.GetModuleId(),
		MapPinTypeGroupID: request.GetMapPinTypeGroupId(),
	})
	if err != nil {
		return nil, err
	}

	err = server.Store.DeleteMapPinTypeGroup(ctx, request.GetMapPinTypeGroupId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func validateDeleteMapPinTypeGroup(req *pb.DeleteMapPinTypeGroupRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUniversalId(req.GetModuleId()); err != nil {
		violations = append(violations, e.FieldViolation("module_id", err))
	}

	if err := validator.ValidateUniversalId(req.GetMapPinTypeGroupId()); err != nil {
		violations = append(violations, e.FieldViolation("map_pin_type_group_id", err))
	}

	return violations
}
