package entities

import (
	"context"
	"github.com/the-medo/talebound-backend/api/servicecore"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *ServiceEntities) DeleteEntityGroupContent(ctx context.Context, request *pb.DeleteEntityGroupContentRequest) (*emptypb.Empty, error) {
	violations := validateDeleteEntityGroupContent(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	err := server.CheckEntityGroupAccess(ctx, request.GetEntityGroupId(), &servicecore.ModulePermission{
		NeedsMenuPermission: true,
	})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to delete entity group content: %v", err)
	}

	arg := db.DeleteEntityGroupContentParams{
		ID:         request.GetContentId(),
		DeleteType: db.DeleteEntityGroupContentActionUnknown,
	}

	err = server.Store.DeleteEntityGroupContent(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func validateDeleteEntityGroupContent(req *pb.DeleteEntityGroupContentRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUniversalId(req.GetEntityGroupId()); err != nil {
		violations = append(violations, e.FieldViolation("entity_group_id", err))
	}
	if err := validator.ValidateUniversalId(req.GetContentId()); err != nil {
		violations = append(violations, e.FieldViolation("content_id", err))
	}
	return violations
}
