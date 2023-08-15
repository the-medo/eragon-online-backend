package api

import (
	"context"
	"database/sql"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *Server) DeleteAvailableWorldTag(ctx context.Context, req *pb.DeleteAvailableWorldTagRequest) (*emptypb.Empty, error) {
	violations := validateDeleteAvailableWorldTagRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	err := server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "can not delete tag - you are not admin: %v", err)
	}

	arg := db.DeleteWorldTagParams{
		TagID: sql.NullInt32{
			Int32: req.GetTagId(),
			Valid: true,
		},
	}

	err = server.store.DeleteWorldTag(ctx, arg)
	if err != nil {
		return nil, err
	}

	err = server.store.DeleteWorldTagAvailable(ctx, req.GetTagId())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func validateDeleteAvailableWorldTagRequest(req *pb.DeleteAvailableWorldTagRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateTagId(req.GetTagId()); err != nil {
		violations = append(violations, FieldViolation("tag_id", err))
	}

	return violations
}
