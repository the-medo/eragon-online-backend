package api

import (
	"context"
	"github.com/the-medo/talebound-backend/api/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) AddWorldTag(ctx context.Context, req *pb.AddWorldTagRequest) (*pb.Tag, error) {
	violations := validateAddWorldTagRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	_, err := server.CheckWorldAdmin(ctx, req.GetWorldId(), false)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to add world tag: %v", err)
	}

	arg := db.CreateWorldTagParams{
		WorldID: req.GetWorldId(),
		TagID:   req.GetTagId(),
	}

	_, err = server.store.CreateWorldTag(ctx, arg)
	if err != nil {
		return nil, err
	}

	tag, err := server.store.GetWorldTagAvailable(ctx, req.GetTagId())
	if err != nil {
		return nil, err
	}

	rsp := converters.ConvertViewTagToTag(tag)

	return rsp, nil
}

func validateAddWorldTagRequest(req *pb.AddWorldTagRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateWorldId(req.GetWorldId()); err != nil {
		violations = append(violations, FieldViolation("world_id", err))
	}

	if err := validator.ValidateTagId(req.GetTagId()); err != nil {
		violations = append(violations, FieldViolation("tag_id", err))
	}

	return violations
}
