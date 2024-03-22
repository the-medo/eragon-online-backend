package entities

import (
	"context"
	"github.com/the-medo/talebound-backend/api/servicecore"
	"github.com/the-medo/talebound-backend/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *ServiceEntities) DeleteEntityGroup(ctx context.Context, request *pb.DeleteEntityGroupRequest) (*emptypb.Empty, error) {
	violations := validateDeleteEntityGroup(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	err := server.CheckEntityGroupAccess(ctx, request.GetEntityGroupId(), &servicecore.ModulePermission{
		NeedsMenuPermission: true,
	})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to delete entity group: %v", err)
	}

	arg := db.DeleteEntityGroupParams{
		ID:         request.GetEntityGroupId(),
		DeleteType: converters.ConvertDeleteEntityGroupContentActionToDB(request.GetDeleteType()),
	}

	err = server.Store.DeleteEntityGroup(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func validateDeleteEntityGroup(req *pb.DeleteEntityGroupRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUniversalId(req.GetEntityGroupId()); err != nil {
		violations = append(violations, e.FieldViolation("entity_group_id", err))
	}
	if err := validator.ValidateDeleteEntityGroupContentAction(req.GetDeleteType()); err != nil {
		violations = append(violations, e.FieldViolation("delete_type", err))
	}
	return violations
}
