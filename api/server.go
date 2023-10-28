package api

import (
	"fmt"
	"github.com/the-medo/talebound-backend/apiservices/entities"
	"github.com/the-medo/talebound-backend/apiservices/locations"
	"github.com/the-medo/talebound-backend/apiservices/maps"
	"github.com/the-medo/talebound-backend/apiservices/modules"
	"github.com/the-medo/talebound-backend/apiservices/srv"
	"github.com/the-medo/talebound-backend/apiservices/tags"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/token"
	"github.com/the-medo/talebound-backend/util"
	"github.com/the-medo/talebound-backend/worker"
)

// Server serves gRPC requests
type Server struct {
	*locations.ServiceLocations
	*maps.ServiceMaps
	*modules.ServiceModules
	*entities.ServiceEntities
	*tags.ServiceTags
	pb.UnimplementedChatServer
	pb.UnimplementedEvaluationsServer
	pb.UnimplementedImagesServer
	pb.UnimplementedMenusServer
	pb.UnimplementedPostTypesServer
	pb.UnimplementedPostsServer
	pb.UnimplementedUsersServer
	pb.UnimplementedVerifyServer
	pb.UnimplementedWorldsServer
	Config          util.Config
	Store           db.Store
	TaskDistributor worker.TaskDistributor
	TokenMaker      token.Maker
}

func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	serverCore := srv.NewServiceCore(config, store, taskDistributor, tokenMaker)

	server := &Server{
		ServiceLocations: locations.NewLocationsService(serverCore),
		ServiceMaps:      maps.NewMapsService(serverCore),
		ServiceModules:   modules.NewModulesService(serverCore),
		ServiceEntities:  entities.NewEntitiesService(serverCore),
		ServiceTags:      tags.NewTagsService(serverCore),
		Config:           serverCore.Config,
		Store:            serverCore.Store,
		TaskDistributor:  serverCore.TaskDistributor,
		TokenMaker:       serverCore.TokenMaker,
	}

	return server, nil
}
