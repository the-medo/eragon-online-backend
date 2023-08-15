package api

import (
	"context"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *Server) DeleteWorldAdmin(ctx context.Context, req *pb.DeleteWorldAdminRequest) (*emptypb.Empty, error) {
	violations := validateDeleteWorldAdmin(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUserCookie(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	//user can remove himself from world admin even if he is not super admin
	_, err = server.CheckWorldAdmin(ctx, req.GetWorldId(), req.GetUserId() != authPayload.UserId)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to delete world admin: %v", err)
	}

	arg := db.DeleteWorldAdminParams{
		WorldID: req.GetWorldId(),
		UserID:  req.GetUserId(),
	}

	err = server.store.DeleteWorldAdmin(ctx, arg)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func validateDeleteWorldAdmin(req *pb.DeleteWorldAdminRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateUserId(req.GetUserId()); err != nil {
		violations = append(violations, FieldViolation("user_id", err))
	}

	if err := validator.ValidateWorldId(req.GetWorldId()); err != nil {
		violations = append(violations, FieldViolation("world_id", err))
	}

	return violations
}
