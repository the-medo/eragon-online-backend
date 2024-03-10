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

func (server *ServiceEntities) UpdateEntityGroupContent(ctx context.Context, req *pb.UpdateEntityGroupContentRequest) (*pb.EntityGroupContent, error) {
	violations := validateUpdateEntityGroupContent(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	if req.ContentEntityId == nil && req.ContentEntityGroupId == nil && req.Position == nil && req.NewEntityGroupId == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nothing to update, all update fields are empty")
	}

	err := server.CheckEntityGroupAccess(ctx, req.GetEntityGroupId(), &servicecore.ModulePermission{
		NeedsMenuPermission: true,
	})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to update entity group content: %v", err)
	}

	if req.Position != nil || req.NewEntityGroupId != nil {
		err = server.Store.EntityGroupContentChangePositions(ctx, db.EntityGroupContentChangePositionsParams{
			ID: req.GetContentId(),
			NewEntityGroupID: sql.NullInt32{
				Int32: req.GetNewEntityGroupId(),
				Valid: req.NewEntityGroupId != nil,
			},
			NewPosition: sql.NullInt32{
				Int32: req.GetPosition(),
				Valid: req.Position != nil,
			},
		})

		if err != nil {
			return nil, err
		}
	}

	arg := db.UpdateEntityGroupContentParams{
		ID: req.GetContentId(),
		ContentEntityID: sql.NullInt32{
			Int32: req.GetContentEntityId(),
			Valid: req.ContentEntityId != nil,
		},
		ContentEntityGroupID: sql.NullInt32{
			Int32: req.GetContentEntityGroupId(),
			Valid: req.ContentEntityGroupId != nil,
		},
	}

	updatedContent, err := server.Store.UpdateEntityGroupContent(ctx, arg)
	if err != nil {
		return nil, err
	}

	rsp := converters.ConvertEntityGroupContent(updatedContent) // Assuming you have a converter for EntityGroupContent

	return rsp, nil
}

func validateUpdateEntityGroupContent(req *pb.UpdateEntityGroupContentRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateUniversalId(req.GetEntityGroupId()); err != nil {
		violations = append(violations, e.FieldViolation("entity_group_id", err))
	}

	if err := validator.ValidateUniversalId(req.GetContentId()); err != nil {
		violations = append(violations, e.FieldViolation("content_id", err))
	}

	if req.NewEntityGroupId != nil {
		if err := validator.ValidateUniversalId(req.GetNewEntityGroupId()); err != nil {
			violations = append(violations, e.FieldViolation("new_entity_group_id", err))
		}
	}

	if req.Position != nil {
		if err := validator.ValidateUniversalId(req.GetPosition()); err != nil {
			violations = append(violations, e.FieldViolation("position", err))
		}
	}

	if req.ContentEntityId != nil {
		if err := validator.ValidateUniversalId(req.GetContentEntityId()); err != nil {
			violations = append(violations, e.FieldViolation("content_entity_id", err))
		}
	}

	if req.ContentEntityGroupId != nil {
		if err := validator.ValidateUniversalId(req.GetContentEntityGroupId()); err != nil {
			violations = append(violations, e.FieldViolation("content_entity_group_id", err))
		}
	}

	return violations
}
