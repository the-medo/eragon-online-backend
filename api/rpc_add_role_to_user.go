package api

import (
	"context"
	"fmt"
	"github.com/the-medo/talebound-backend/api/e"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) AddRoleToUser(ctx context.Context, req *pb.AddRoleToUserRequest) (*pb.AddRoleToUserResponse, error) {
	violations := validateAddRoleToUser(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	err := server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to add role to user: %v", err)
	}

	println(req.GetUserId(), req.GetRoleId())
	newRole, err := server.Store.AddUserRole(ctx, db.AddUserRoleParams{
		UserID: req.GetUserId(),
		RoleID: req.GetRoleId(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add role to user: %v", err)
	}

	rsp := &pb.AddRoleToUserResponse{
		Success: true,
		Message: fmt.Sprintf("Role %s added to user successfully.", pb.RoleType_name[newRole.RoleID]),
	}

	return rsp, nil
}

func validateAddRoleToUser(req *pb.AddRoleToUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateUserId(req.GetUserId()); err != nil {
		violations = append(violations, e.FieldViolation("user_id", err))
	}
	if err := validator.ValidateRoleId(req.GetRoleId()); err != nil {
		violations = append(violations, e.FieldViolation("role_id", err))
	}

	return violations
}
