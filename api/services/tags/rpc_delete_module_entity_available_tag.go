package tags

import (
	"context"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *ServiceTags) DeleteModuleEntityAvailableTag(ctx context.Context, req *pb.DeleteModuleEntityAvailableTagRequest) (*emptypb.Empty, error) {

	violations := validateDeleteModuleEntityAvailableTagRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	err := server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "can not update tag - you are not admin: %v", err)
	}

	err = server.Store.DeleteModuleEntityTagAvailable(ctx, req.GetTagId())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func validateDeleteModuleEntityAvailableTagRequest(req *pb.DeleteModuleEntityAvailableTagRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateTagId(req.GetTagId()); err != nil {
		violations = append(violations, e.FieldViolation("tag_id", err))
	}

	return violations
}
