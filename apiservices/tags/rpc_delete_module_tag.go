package tags

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/api/e"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *ServiceTags) DeleteModuleTag(ctx context.Context, req *pb.DeleteModuleTagRequest) (*emptypb.Empty, error) {

	violations := validateDeleteModuleTagRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	err := server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "can not update tag - you are not admin: %v", err)
	}

	err = server.Store.DeleteModuleTag(ctx, db.DeleteModuleTagParams{
		ModuleID: sql.NullInt32{
			Int32: req.GetModuleId(),
			Valid: true,
		},
		TagID: sql.NullInt32{
			Int32: req.GetTagId(),
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func validateDeleteModuleTagRequest(req *pb.DeleteModuleTagRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateUniversalId(req.GetModuleId()); err != nil {
		violations = append(violations, e.FieldViolation("module_id", err))
	}

	if err := validator.ValidateTagId(req.GetTagId()); err != nil {
		violations = append(violations, e.FieldViolation("tag_id", err))
	}

	return violations
}
