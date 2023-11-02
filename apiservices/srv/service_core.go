package srv

import (
	"github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/token"
	"github.com/the-medo/talebound-backend/util"
	"github.com/the-medo/talebound-backend/worker"
)

const GrpcCookieHeader = "grpc-gateway-cookie"
const CookieName = "access_token"

type ServiceCore struct {
	Config          util.Config
	Store           db.Store
	TaskDistributor worker.TaskDistributor
	TokenMaker      token.Maker
}

type ModulePermission struct {
	NeedsSuperAdmin       bool
	NeedsMenuPermission   bool
	NeedsEntityPermission *[]db.EntityType
}

func NewServiceCore(config util.Config, store db.Store, taskDistributor worker.TaskDistributor, tokenMaker token.Maker) *ServiceCore {
	return &ServiceCore{
		Config:          config,
		Store:           store,
		TaskDistributor: taskDistributor,
		TokenMaker:      tokenMaker,
	}
}
