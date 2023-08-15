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

func (server *Server) RemoveWorldTag(ctx context.Context, req *pb.RemoveWorldTagRequest) (*emptypb.Empty, error) {
	violations := validateRemoveWorldTag(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	_, err := server.CheckWorldAdmin(ctx, req.GetWorldId(), false)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to remove world tag: %v", err)
	}

	arg := db.DeleteWorldTagParams{
		TagID: sql.NullInt32{
			Int32: req.GetTagId(),
			Valid: true,
		},
		WorldID: sql.NullInt32{
			Int32: req.GetWorldId(),
			Valid: true,
		},
	}

	err = server.store.DeleteWorldTag(ctx, arg)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func validateRemoveWorldTag(req *pb.RemoveWorldTagRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateWorldId(req.GetWorldId()); err != nil {
		violations = append(violations, FieldViolation("world_id", err))
	}

	if err := validator.ValidateTagId(req.GetTagId()); err != nil {
		violations = append(violations, FieldViolation("tag_id", err))
	}

	return violations
}
