package entities

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/api/e"
	"github.com/the-medo/talebound-backend/apiservices/servicecore"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *ServiceEntities) UpdateEntityGroupContent(ctx context.Context, request *pb.UpdateEntityGroupContentRequest) (*pb.EntityGroupContent, error) {
	violations := validateUpdateEntityGroupContent(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	err := server.CheckEntityGroupAccess(ctx, request.GetEntityGroupId(), &servicecore.ModulePermission{
		NeedsMenuPermission: true,
	})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to update entity group content: %v", err)
	}

	arg := db.UpdateEntityGroupContentParams{
		ID: request.GetContentId(),
		NewEntityGroupID: sql.NullInt32{
			Int32: request.GetNewEntityGroupId(),
			Valid: request.NewEntityGroupId != nil,
		},
		ContentEntityID: sql.NullInt32{
			Int32: request.GetContentEntityId(),
			Valid: request.ContentEntityId != nil,
		},
		ContentEntityGroupID: sql.NullInt32{
			Int32: request.GetContentEntityGroupId(),
			Valid: request.ContentEntityGroupId != nil,
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

	return violations
}
