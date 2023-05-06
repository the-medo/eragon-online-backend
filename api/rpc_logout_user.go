package api

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *Server) LogoutUser(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	md := metadata.Pairs(
		"X-Access-Token", "remove",
		"X-Access-Token-Expires-At", "remove",
		"X-Refresh-Token", "remove",
		"X-Refresh-Token-Expires-At", "remove",
	)

	err := grpc.SendHeader(ctx, md)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
