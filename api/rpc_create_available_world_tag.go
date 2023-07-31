package api

import (
	"context"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateAvailableWorldTag(ctx context.Context, request *pb.CreateAvailableWorldTagRequest) (*pb.Tag, error) {
	err := server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "can not create new tag - you are not admin: %v", err)
	}

	tag, err := server.store.CreateWorldTagAvailable(ctx, request.GetTag())
	if err != nil {
		return nil, err
	}

	rsp := converters.ConvertTag(tag)

	return rsp, nil
}
