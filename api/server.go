package api

import (
	"fmt"
	"github.com/the-medo/talebound-backend/apiservices/auth"
	"github.com/the-medo/talebound-backend/apiservices/entities"
	"github.com/the-medo/talebound-backend/apiservices/evaluations"
	"github.com/the-medo/talebound-backend/apiservices/images"
	"github.com/the-medo/talebound-backend/apiservices/locations"
	"github.com/the-medo/talebound-backend/apiservices/maps"
	"github.com/the-medo/talebound-backend/apiservices/menus"
	"github.com/the-medo/talebound-backend/apiservices/modules"
	"github.com/the-medo/talebound-backend/apiservices/posts"
	"github.com/the-medo/talebound-backend/apiservices/srv"
	"github.com/the-medo/talebound-backend/apiservices/tags"
	"github.com/the-medo/talebound-backend/apiservices/users"
	"github.com/the-medo/talebound-backend/apiservices/worlds"
	db "github.com/the-medo/talebound-backend/db/sqlc"
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
	*users.ServiceUsers
	*menus.ServiceMenus
	*posts.ServicePosts
	*evaluations.ServiceEvaluations
	*images.ServiceImages
	*auth.ServiceAuth
	*worlds.ServiceWorlds
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
		ServiceLocations:   locations.NewLocationsService(serverCore),
		ServiceMaps:        maps.NewMapsService(serverCore),
		ServiceModules:     modules.NewModulesService(serverCore),
		ServiceEntities:    entities.NewEntitiesService(serverCore),
		ServiceTags:        tags.NewTagsService(serverCore),
		ServiceUsers:       users.NewUsersService(serverCore),
		ServiceMenus:       menus.NewMenusService(serverCore),
		ServicePosts:       posts.NewPostsService(serverCore),
		ServiceEvaluations: evaluations.NewEvaluationsService(serverCore),
		ServiceImages:      images.NewImagesService(serverCore),
		ServiceAuth:        auth.NewAuthService(serverCore),
		ServiceWorlds:      worlds.NewWorldsService(serverCore),
		Config:             serverCore.Config,
		Store:              serverCore.Store,
		TaskDistributor:    serverCore.TaskDistributor,
		TokenMaker:         serverCore.TokenMaker,
	}

	return server, nil
}
