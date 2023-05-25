package api

import (
	"context"
	"database/sql"
	"fmt"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	limit := req.GetLimit()
	if limit == 0 {
		limit = 100
	}

	offset := req.GetOffset()
	if offset < 0 {
		offset = 0
	}

	users, err := server.store.GetUsers(ctx, db.GetUsersParams{
		PageLimit:  limit,
		PageOffset: offset,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get users: %v", err)
	}

	rsp := &pb.GetUsersResponse{
		Users: make([]*pb.User, len(users)),
	}

	for i, user := range users {
		rsp.Users[i] = convertUserRowWithImage(user)
	}

	return rsp, nil
}

func (server *Server) GetUserRoles(ctx context.Context, req *pb.GetUserRolesRequest) (*pb.GetUserRolesResponse, error) {
	roles, err := server.store.GetUserRoles(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user roles: %v", err)
	}

	rsp := &pb.GetUserRolesResponse{
		Role: make([]*pb.Role, len(roles)),
	}

	for i, role := range roles {
		rsp.Role[i] = &pb.Role{
			Id:          role.RoleID,
			Name:        role.RoleName,
			Description: role.RoleDescription,
		}
	}

	return rsp, nil
}

func (server *Server) AddRoleToUser(ctx context.Context, req *pb.AddRoleToUserRequest) (*pb.AddRoleToUserResponse, error) {
	err := server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to add role to user: %v", err)
	}

	println(req.GetUserId(), req.GetRoleId())
	newRole, err := server.store.AddUserRole(ctx, db.AddUserRoleParams{
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

func (server *Server) RemoveRoleFromUser(ctx context.Context, req *pb.RemoveRoleFromUserRequest) (*pb.RemoveRoleFromUserResponse, error) {
	err := server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to remove role from user: %v", err)
	}

	_, err = server.store.HasUserRole(ctx, db.HasUserRoleParams{
		UserID: req.GetUserId(),
		Role:   pb.RoleType_name[req.GetRoleId()],
	})

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user doesn't have that role: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to remove role from user: %v", err)
	}

	err = server.store.RemoveUserRole(ctx, db.RemoveUserRoleParams{
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
