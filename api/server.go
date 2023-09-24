package api

import (
	"fmt"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/token"
	"github.com/the-medo/talebound-backend/util"
	"github.com/the-medo/talebound-backend/worker"
)

//err = pb.RegisterChatHandlerServer(ctx, grpcMux, server)
//err = pb.RegisterEvaluationsHandlerServer(ctx, grpcMux, server)
//err = pb.RegisterImagesHandlerServer(ctx, grpcMux, server)
//err = pb.RegisterMapsHandlerServer(ctx, grpcMux, server)
//err = pb.RegisterMenusHandlerServer(ctx, grpcMux, server)
//err = pb.RegisterPostTypesHandlerServer(ctx, grpcMux, server)
//err = pb.RegisterPostsHandlerServer(ctx, grpcMux, server)
//err = pb.RegisterTagsHandlerServer(ctx, grpcMux, server)
//err = pb.RegisterUsersHandlerServer(ctx, grpcMux, server)
//err = pb.RegisterVerifyHandlerServer(ctx, grpcMux, server)
//err = pb.RegisterWorldsHandlerServer(ctx, grpcMux, server)

// Server serves gRPC requests
type Server struct {
	pb.UnimplementedChatServer
	pb.UnimplementedEvaluationsServer
	pb.UnimplementedImagesServer
	pb.UnimplementedLocationsServer
	pb.UnimplementedMapsServer
	pb.UnimplementedMenusServer
	pb.UnimplementedPostTypesServer
	pb.UnimplementedPostsServer
	pb.UnimplementedTagsServer
	pb.UnimplementedUsersServer
	pb.UnimplementedVerifyServer
	pb.UnimplementedWorldsServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
