package api

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/the-medo/talebound-backend/api/e"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) RemoveRoleFromUser(ctx context.Context, req *pb.RemoveRoleFromUserRequest) (*pb.RemoveRoleFromUserResponse, error) {
	violations := validateRemoveRoleFromUser(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	err := server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to remove role from user: %v", err)
	}

	_, err = server.Store.HasUserRole(ctx, db.HasUserRoleParams{
		UserID: req.GetUserId(),
		Role:   pb.RoleType_name[req.GetRoleId()],
	})

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user doesn't have that role: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to remove role from user: %v", err)
	}

	err = server.Store.RemoveUserRole(ctx, db.RemoveUserRoleParams{
		UserID: req.GetUserId(),
		RoleID: req.GetRoleId(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to remove role from user: %v", err)
	}

	rsp := &pb.RemoveRoleFromUserResponse{
		Success: true,
		Message: fmt.Sprintf("Role %s removed from user successfully.", pb.RoleType_name[req.GetRoleId()]),
	}

	return rsp, nil
}

func validateRemoveRoleFromUser(req *pb.RemoveRoleFromUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateUserId(req.GetUserId()); err != nil {
		violations = append(violations, e.FieldViolation("user_id", err))
	}
	if err := validator.ValidateRoleId(req.GetRoleId()); err != nil {
		violations = append(violations, e.FieldViolation("role_id", err))
	}

	return violations
}
