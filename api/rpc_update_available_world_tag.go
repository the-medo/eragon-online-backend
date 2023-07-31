package api

import (
	"context"
	"github.com/the-medo/talebound-backend/api/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateAvailableWorldTag(ctx context.Context, request *pb.UpdateAvailableWorldTagRequest) (*pb.Tag, error) {
	err := server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "can not update tag - you are not admin: %v", err)
	}

	arg := db.UpdateWorldTagAvailableParams{
		ID:  request.GetTagId(),
		Tag: request.GetNewTag(),
	}

	tag, err := server.store.UpdateWorldTagAvailable(ctx, arg)
	if err != nil {
		return nil, err
	}

	rsp := converters.ConvertTag(tag)

	return rsp, nil
}
