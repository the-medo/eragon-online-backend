package api

import (
	"context"
	"github.com/the-medo/talebound-backend/api/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) AddWorldTag(ctx context.Context, req *pb.AddWorldTagRequest) (*pb.Tag, error) {

	err := server.CheckWorldAdmin(ctx, req.GetWorldId(), false)
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

	rsp := converters.ConvertTag(tag)

	return rsp, nil
}
