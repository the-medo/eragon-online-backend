package api

import (
	"context"
	"database/sql"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *Server) RemoveWorldTag(ctx context.Context, req *pb.RemoveWorldTagRequest) (*emptypb.Empty, error) {
	err := server.CheckWorldAdmin(ctx, req.GetWorldId(), false)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to remove world tag: %v", err)
	}

	arg := db.DeleteWorldTagParams{
		TagID: sql.NullInt32{
			Int32: req.GetTagId(),
			Valid: true,
		},
		WorldID: sql.NullInt32{
			Int32: req.GetWorldId(),
			Valid: true,
		},
	}

	err = server.store.DeleteWorldTag(ctx, arg)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
