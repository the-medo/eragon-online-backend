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

func (server *ServiceEntities) CreateEntityGroup(ctx context.Context, request *pb.CreateEntityGroupRequest) (*pb.CreateEntityGroupResponse, error) {
	violations := validateCreateEntityGroup(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	err := server.CheckEntityGroupAccess(ctx, request.GetParentEntityGroupId(), &servicecore.ModulePermission{
		NeedsMenuPermission: true,
	})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to create entity group: %v", err)
	}

	arg := db.CreateEntityGroupParams{
		Name: sql.NullString{
			String: request.GetName(),
			Valid:  request.Name != nil,
		},
		Description: sql.NullString{
			String: request.GetDescription(),
			Valid:  request.Description != nil,
		},
		Style: sql.NullString{
			String: converters.ConvertEntityGroupStyleToDB(request.GetStyle()),
			Valid:  request.Style != nil,
		},
		Direction: sql.NullString{
			String: converters.ConvertEntityGroupDirectionToDB(request.GetDirection()),
			Valid:  request.Direction != nil,
		},
	}

	newEntityGroup, err := server.Store.CreateEntityGroup(ctx, arg)
	if err != nil {
		return nil, err
	}

	arg2 := db.CreateEntityGroupContentParams{
		EntityGroupID: request.GetParentEntityGroupId(),
		ContentEntityGroupID: sql.NullInt32{
			Int32: newEntityGroup.ID,
			Valid: true,
		},
		Position: sql.NullInt32{
			Int32: request.GetPosition(),
			Valid: request.Position != nil,
		},
	}

	newEntityGroupContent, err := server.Store.CreateEntityGroupContent(ctx, arg2)
	if err != nil {
		return nil, err
	}

	rsp := &pb.CreateEntityGroupResponse{
		EntityGroup:        converters.ConvertEntityGroup(newEntityGroup),
		EntityGroupContent: converters.ConvertEntityGroupContent(newEntityGroupContent),
	}

	return rsp, nil
}

func validateCreateEntityGroup(req *pb.CreateEntityGroupRequest) (violations []*errdetails.BadRequest_FieldViolation) {
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

	if req.Position != nil {
		if err := validator.ValidateInt(req.GetPosition(), 1); err != nil {
			violations = append(violations, e.FieldViolation("position", err))
		}
	}

	if req.Style != nil {
		if err := validator.ValidateEntityGroupStyle(converters.ConvertEntityGroupStyleToDB(req.GetStyle())); err != nil {
			violations = append(violations, e.FieldViolation("style", err))
		}
	}

	if req.Direction != nil {
		if err := validator.ValidateEntityGroupDirection(converters.ConvertEntityGroupDirectionToDB(req.GetDirection())); err != nil {
			violations = append(violations, e.FieldViolation("direction", err))
		}
	}

	return violations
}
