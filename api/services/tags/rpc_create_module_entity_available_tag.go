package tags

import (
	"context"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *ServiceTags) CreateModuleEntityAvailableTag(ctx context.Context, req *pb.CreateModuleEntityAvailableTagRequest) (*pb.EntityTagAvailable, error) {

	violations := validateCreateModuleEntityAvailableTagRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	err := server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "can not create new tag - you are not admin: %v", err)
	}

	tag, err := server.Store.CreateModuleEntityTagAvailable(ctx, db.CreateModuleEntityTagAvailableParams{
		ModuleID: req.GetModuleId(),
		Tag:      req.GetTag(),
	})
	if err != nil {
		return nil, err
	}

	rsp := &pb.EntityTagAvailable{
		Id:       tag.ID,
		Tag:      tag.Tag,
		ModuleId: tag.ModuleID,
	}

	return rsp, nil
}

func validateCreateModuleEntityAvailableTagRequest(req *pb.CreateModuleEntityAvailableTagRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateTag(req.GetTag()); err != nil {
		violations = append(violations, e.FieldViolation("tag", err))
	}

	return violations
}
