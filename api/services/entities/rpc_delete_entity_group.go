package entities

import (
	"context"
	"github.com/the-medo/talebound-backend/api/servicecore"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *ServiceEntities) DeleteEntityGroup(ctx context.Context, request *pb.DeleteEntityGroupRequest) (*emptypb.Empty, error) {
	err := server.CheckEntityGroupAccess(ctx, request.GetEntityGroupId(), &servicecore.ModulePermission{
		NeedsMenuPermission: true,
	})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to delete entity group: %v", err)
	}

	err = server.Store.DeleteEntityGroup(ctx, request.GetEntityGroupId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
