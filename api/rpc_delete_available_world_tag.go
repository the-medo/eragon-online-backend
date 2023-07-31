package api

import (
	"context"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *Server) DeleteAvailableWorldTag(ctx context.Context, request *pb.DeleteAvailableWorldTagRequest) (*emptypb.Empty, error) {
	err := server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "can not delete tag - you are not admin: %v", err)
	}

	err = server.store.DeleteWorldTagAvailable(ctx, request.GetTagId())
	if err != nil {
		return nil, err
	}

	return nil, nil
}
