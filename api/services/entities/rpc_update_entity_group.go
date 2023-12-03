package entities

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

func (server *ServiceEntities) UpdateEntityGroup(ctx context.Context, request *pb.UpdateEntityGroupRequest) (*pb.EntityGroup, error) {
	violations := validateUpdateEntityGroup(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	err := server.CheckEntityGroupAccess(ctx, request.GetEntityGroupId(), &servicecore.ModulePermission{
		NeedsMenuPermission: true,
	})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to update entity group: %v", err)
	}

	arg := db.UpdateEntityGroupParams{
		ID: request.GetEntityGroupId(),
		Name: sql.NullString{
			String: request.GetName(),
			Valid:  request.Name != nil,
		},
		Description: sql.NullString{
			String: request.GetDescription(),
			Valid:  request.Description != nil,
		},
		Style: sql.NullString{
			String: request.GetStyle(),
			Valid:  request.Style != nil,
		},
		Direction: sql.NullString{
			String: request.GetDirection(),
			Valid:  request.Direction != nil,
		},
	}

	updatedEntityGroup, err := server.Store.UpdateEntityGroup(ctx, arg)
	if err != nil {
		return nil, err
	}

	rsp := converters.ConvertEntityGroup(updatedEntityGroup) // Assuming you have a converter for EntityGroup

	return rsp, nil
}

func validateUpdateEntityGroup(req *pb.UpdateEntityGroupRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if req.Name != nil {
		if err := validator.ValidateUniversalName(req.GetName()); err != nil {
			violations = append(violations, e.FieldViolation("name", err))
		}
	}

	if req.Description != nil {
		if err := validator.ValidateUniversalDescription(req.GetDescription()); err != nil {
			violations = append(violations, e.FieldViolation("description", err))
		}
	}

	if req.Style != nil {
		if err := validator.ValidateEntityGroupStyle(req.GetStyle()); err != nil {
			violations = append(violations, e.FieldViolation("style", err))
		}
	}

	if req.Direction != nil {
		if err := validator.ValidateEntityGroupDirection(req.GetDirection()); err != nil {
			violations = append(violations, e.FieldViolation("direction", err))
		}
	}

	return violations
}
