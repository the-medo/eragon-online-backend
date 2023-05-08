package api

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *Server) LogoutUser(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	md := metadata.Pairs(
		"X-Access-Token", "null",
		"X-Access-Token-Expires-At", "null",
		"X-Refresh-Token", "null",
		"X-Refresh-Token-Expires-At", "null",
	)

	err := grpc.SendHeader(ctx, md)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
