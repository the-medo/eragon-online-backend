package srv

import (
	"fmt"
	"github.com/the-medo/talebound-backend/api/servicecore"
	"github.com/the-medo/talebound-backend/api/services/auth"
	"github.com/the-medo/talebound-backend/api/services/characters"
	"github.com/the-medo/talebound-backend/api/services/entities"
	"github.com/the-medo/talebound-backend/api/services/evaluations"
	"github.com/the-medo/talebound-backend/api/services/fetcher"
	"github.com/the-medo/talebound-backend/api/services/images"
	"github.com/the-medo/talebound-backend/api/services/locations"
	"github.com/the-medo/talebound-backend/api/services/maps"
	"github.com/the-medo/talebound-backend/api/services/menus"
	"github.com/the-medo/talebound-backend/api/services/modules"
	"github.com/the-medo/talebound-backend/api/services/posts"
	"github.com/the-medo/talebound-backend/api/services/quests"
	"github.com/the-medo/talebound-backend/api/services/systems"
	"github.com/the-medo/talebound-backend/api/services/tags"
	"github.com/the-medo/talebound-backend/api/services/users"
	"github.com/the-medo/talebound-backend/api/services/worlds"
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
	*systems.ServiceSystems
	*characters.ServiceCharacters
	*quests.ServiceQuests
	*fetcher.ServiceFetcher
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

	serverCore := servicecore.NewServiceCore(config, store, taskDistributor, tokenMaker)

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
		ServiceSystems:     systems.NewSystemsService(serverCore),
		ServiceCharacters:  characters.NewCharactersService(serverCore),
		ServiceQuests:      quests.NewQuestsService(serverCore),
		ServiceFetcher:     fetcher.NewFetcherService(serverCore),
		Config:             serverCore.Config,
		Store:              serverCore.Store,
		TaskDistributor:    serverCore.TaskDistributor,
		TokenMaker:         serverCore.TokenMaker,
	}

	return server, nil
}
