package entities

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/api/e"
	"github.com/the-medo/talebound-backend/apiservices/servicecore"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
		ContentEntityID: sql.NullInt32{
			Int32: request.GetContentEntityId(),
			Valid: request.ContentEntityId != nil,
		},
		ContentEntityGroupID: sql.NullInt32{
			Int32: request.GetContentEntityGroupId(),
			Valid: request.ContentEntityGroupId != nil,
		},
	}

	newContent, err := server.Store.CreateEntityGroupContent(ctx, arg)
	if err != nil {
		return nil, err
	}

	rsp := converters.ConvertEntityGroupContent(newContent) // Assuming you have a converter for EntityGroupContent

	return rsp, nil
}

func validateCreateEntityGroupContent(req *pb.CreateEntityGroupContentRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if req.ContentEntityId == nil && req.ContentEntityGroupId == nil {
		violations = append(violations, e.FieldViolation("content_entity_id/content_entity_group_id", fmt.Errorf("at least one of the content fields must be provided")))
	}

	return violations
}
