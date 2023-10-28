package tags

import (
	"context"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/api/e"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *ServiceTags) CreateModuleTypeAvailableTag(ctx context.Context, req *pb.CreateModuleTypeAvailableTagRequest) (*pb.ViewTag, error) {

	violations := validateCreateModuleTypeAvailableTagRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	err := server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "can not create new tag - you are not admin: %v", err)
	}

	tag, err := server.Store.CreateModuleTypeTagAvailable(ctx, db.CreateModuleTypeTagAvailableParams{
		ModuleType: converters.ConvertModuleTypeToDB(req.GetModuleType()),
		Tag:        req.GetTag(),
	})
	if err != nil {
		return nil, err
	}

	rsp := &pb.ViewTag{
		Id:         tag.ID,
		Tag:        tag.Tag,
		ModuleType: converters.ConvertModuleTypeToPB(tag.ModuleType),
		Count:      0,
	}

	return rsp, nil
}

func validateCreateModuleTypeAvailableTagRequest(req *pb.CreateModuleTypeAvailableTagRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateTag(req.GetTag()); err != nil {
		violations = append(violations, e.FieldViolation("tag", err))
	}

	return violations
}
