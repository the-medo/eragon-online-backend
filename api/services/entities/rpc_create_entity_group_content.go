package entities

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

// CreateEntityGroupContent adds entity_group_content based on request
//  1. content_entity_group_id => in case of adding group, only this field is necessary
//  2. content_entity_id => if entityId is known, only this field is necessary
//  3. entity_type + entity_id_of_type => combination of both is needed to create new entity (if necessary) and insert it as content
func (server *ServiceEntities) CreateEntityGroupContent(ctx context.Context, request *pb.CreateEntityGroupContentRequest) (*pb.EntityGroupContent, error) {
	violations := validateCreateEntityGroupContent(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	err := server.CheckEntityGroupAccess(ctx, request.GetEntityGroupId(), &servicecore.ModulePermission{
		NeedsMenuPermission: true,
	})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to create entity group content: %v", err)
	}

	arg := db.CreateEntityGroupContentParams{
		EntityGroupID: request.GetEntityGroupId(),
		Position: sql.NullInt32{
			Int32: request.GetPosition(),
			Valid: request.Position != nil,
		},
	}

	if request.ContentEntityGroupId != nil {
		arg.ContentEntityGroupID = sql.NullInt32{
			Int32: request.GetContentEntityGroupId(),
			Valid: true,
		}
	} else if request.ContentEntityId != nil {
		arg.ContentEntityID = sql.NullInt32{
			Int32: request.GetContentEntityId(),
			Valid: true,
		}
	} else if request.EntityType != nil && request.EntityIdOfType != nil {
		// 1. look for entity, if exists and get ENTITY_ID
		entityArg := sql.NullInt32{
			Int32: request.GetEntityIdOfType(),
			Valid: true,
		}
		var err error
		var entity db.Entity

		convertedEntityType := converters.ConvertEntityTypeToDB(request.GetEntityType())

		if convertedEntityType == db.EntityTypePost {
			entity, err = server.Store.GetEntityByPostId(ctx, entityArg)
		} else if convertedEntityType == db.EntityTypeMap {
			entity, err = server.Store.GetEntityByMapId(ctx, entityArg)
		} else if convertedEntityType == db.EntityTypeLocation {
			entity, err = server.Store.GetEntityByLocationId(ctx, entityArg)
		} else if convertedEntityType == db.EntityTypeImage {
			entity, err = server.Store.GetEntityByImageId(ctx, entityArg)
		}

		// 2. create new entity, if it doesn't exist
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				//no entity found, we have to create it

				moduleId, err := server.Store.GetModuleIdOfEntityGroup(ctx, request.GetEntityGroupId())
				if err != nil {
					return nil, fmt.Errorf("failed to get module id of entity group: %w", err)
				}

				entity, err = server.Store.CreateEntity(ctx, db.CreateEntityParams{
					Type:     convertedEntityType,
					ModuleID: moduleId,
					PostID: sql.NullInt32{
						Int32: request.GetEntityIdOfType(),
						Valid: convertedEntityType == db.EntityTypePost,
					},
					MapID: sql.NullInt32{
						Int32: request.GetEntityIdOfType(),
						Valid: convertedEntityType == db.EntityTypeMap,
					},
					LocationID: sql.NullInt32{
						Int32: request.GetEntityIdOfType(),
						Valid: convertedEntityType == db.EntityTypeLocation,
					},
					ImageID: sql.NullInt32{
						Int32: request.GetEntityIdOfType(),
						Valid: convertedEntityType == db.EntityTypeImage,
					},
				})

			} else {
				return nil, fmt.Errorf("failed to get entity: %w", err)
			}
		}

		// 3. update arg
		arg.ContentEntityID = sql.NullInt32{
			Int32: entity.ID,
			Valid: true,
		}
	} else {
		//this shouldn't happen
		return nil, fmt.Errorf("invalid state in rpc_create_entity_group_content")
	}

	newContent, err := server.Store.CreateEntityGroupContent(ctx, arg)
	if err != nil {
		return nil, err
	}

	rsp := converters.ConvertEntityGroupContent(newContent) // Assuming you have a converter for EntityGroupContent

	return rsp, nil
}

func validateCreateEntityGroupContent(req *pb.CreateEntityGroupContentRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if req.ContentEntityId == nil && req.ContentEntityGroupId == nil && req.EntityType == nil && req.EntityIdOfType == nil {
		violations = append(violations, e.FieldViolation("content_entity_id/content_entity_group_id/entity_type/entity_id_of_type", fmt.Errorf("one of these must be provided: content_entity_id OR content_entity_group_id OR (entity_type AND entity_id_of_type)")))
	}

	if req.ContentEntityGroupId != nil {
		if req.ContentEntityId != nil || req.EntityType != nil || req.EntityIdOfType != nil {
			violations = append(violations, e.FieldViolation("content_entity_group_id", fmt.Errorf("when content_entity_group_id is provided, fields content_entity_id, entity_type and entity_id_of_type must be empty")))
		} else {
			if err := validator.ValidateUniversalId(req.GetContentEntityGroupId()); err != nil {
				violations = append(violations, e.FieldViolation("content_entity_group_id", err))
			}
		}
	}

	if req.ContentEntityId != nil {
		if req.ContentEntityGroupId != nil || req.EntityType != nil || req.EntityIdOfType != nil {
			violations = append(violations, e.FieldViolation("content_entity_id", fmt.Errorf("when content_entity_id is provided, fields content_entity_group_id, entity_type and entity_id_of_type must be empty")))
		} else {
			if err := validator.ValidateUniversalId(req.GetContentEntityId()); err != nil {
				violations = append(violations, e.FieldViolation("content_entity_id", err))
			}
		}
	}

	if req.EntityType != nil {
		//TODO: add entity_type validator
		if req.EntityIdOfType != nil {
			if err := validator.ValidateUniversalId(req.GetEntityIdOfType()); err != nil {
				violations = append(violations, e.FieldViolation("entity_id_of_type", err))
			}
		} else {
			violations = append(violations, e.FieldViolation("entity_id_of_type", fmt.Errorf("when entity_type is provided, field entity_id_of_type can not be empty")))
		}
	} else if req.EntityIdOfType != nil {
		violations = append(violations, e.FieldViolation("entity_type", fmt.Errorf("when entity_id_of_type is provided, field entity_type can not be empty")))
	}

	return violations
}
