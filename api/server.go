package api

import (
	"github.com/rs/zerolog/log"
	"github.com/the-medo/talebound-backend/apiservices/locations"
	"github.com/the-medo/talebound-backend/apiservices/maps"
	"github.com/the-medo/talebound-backend/apiservices/srv"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/util"
	"github.com/the-medo/talebound-backend/worker"
)

// Server serves gRPC requests
type Server struct {
	locations.ServerLocations
	maps.ServerMaps
	pb.UnimplementedChatServer
	pb.UnimplementedEvaluationsServer
	pb.UnimplementedImagesServer
	pb.UnimplementedMenusServer
	pb.UnimplementedPostTypesServer
	pb.UnimplementedPostsServer
	pb.UnimplementedTagsServer
	pb.UnimplementedUsersServer
	pb.UnimplementedVerifyServer
	pb.UnimplementedWorldsServer
}

func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {

	serverCore, err := srv.NewServerCore(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create server core:")
	}

	server := &Server{
		ServerLocations: locations.NewLocationsServer(serverCore),
		ServerMaps:      maps.NewMapsServer(serverCore),
	}

	return server, nil
}
