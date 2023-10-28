package tags

import (
	"context"
	"github.com/the-medo/talebound-backend/api/e"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *ServiceTags) UpdateModuleTypeAvailableTag(ctx context.Context, req *pb.UpdateModuleTypeAvailableTagRequest) (*pb.ViewTag, error) {

	violations := validateUpdateModuleTypeAvailableTagRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	err := server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "can not update tag - you are not admin: %v", err)
	}

	tag, err := server.Store.UpdateModuleTypeTagAvailable(ctx, db.UpdateModuleTypeTagAvailableParams{
		ID:  req.GetTagId(),
		Tag: req.GetNewTag(),
	})
	if err != nil {
		return nil, err
	}

	rsp := &pb.ViewTag{
		Id:    tag.ID,
		Tag:   tag.Tag,
		Count: 0,
	}

	return rsp, nil
}

func validateUpdateModuleTypeAvailableTagRequest(req *pb.UpdateModuleTypeAvailableTagRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateTag(req.GetNewTag()); err != nil {
		violations = append(violations, e.FieldViolation("new_tag", err))
	}

	return violations
}
