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

func (server *ServiceTags) CreateEntityTag(ctx context.Context, req *pb.CreateEntityTagRequest) (*pb.CreateEntityTagResponse, error) {

	violations := validateCreateEntityTagRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	err := server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "can not create new tag - you are not admin: %v", err)
	}

	tag, err := server.Store.CreateEntityTag(ctx, db.CreateEntityTagParams{
		EntityID: req.GetEntityId(),
		TagID:    req.GetTagId(),
	})
	if err != nil {
		return nil, err
	}

	rsp := &pb.CreateEntityTagResponse{
		EntityId: tag.EntityID,
		TagId:    tag.TagID,
	}

	return rsp, nil
}

func validateCreateEntityTagRequest(req *pb.CreateEntityTagRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateUniversalId(req.GetTagId()); err != nil {
		violations = append(violations, e.FieldViolation("module_id", err))
	}

	if err := validator.ValidateTagId(req.GetTagId()); err != nil {
		violations = append(violations, e.FieldViolation("tag_id", err))
	}

	return violations
}
