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

func (server *ServiceTags) CreateModuleTag(ctx context.Context, req *pb.CreateModuleTagRequest) (*pb.CreateModuleTagResponse, error) {

	violations := validateCreateModuleTagRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	err := server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "can not create new tag - you are not admin: %v", err)
	}

	tag, err := server.Store.CreateModuleTag(ctx, db.CreateModuleTagParams{
		ModuleID: req.GetModuleId(),
		TagID:    req.GetTagId(),
	})
	if err != nil {
		return nil, err
	}

	rsp := &pb.CreateModuleTagResponse{
		ModuleId: tag.ModuleID,
		TagId:    tag.TagID,
	}

	return rsp, nil
}

func validateCreateModuleTagRequest(req *pb.CreateModuleTagRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateUniversalId(req.GetTagId()); err != nil {
		violations = append(violations, e.FieldViolation("module_id", err))
	}

	if err := validator.ValidateTagId(req.GetTagId()); err != nil {
		violations = append(violations, e.FieldViolation("tag_id", err))
	}

	return violations
}
