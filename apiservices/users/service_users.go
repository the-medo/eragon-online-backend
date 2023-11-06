package users

import (
	"github.com/the-medo/talebound-backend/apiservices/servicecore"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/token"
	"github.com/the-medo/talebound-backend/util"
	"github.com/the-medo/talebound-backend/worker"
	"testing"
	"time"
)

type ServiceUsers struct {
	pb.UnimplementedUsersServer
	*servicecore.ServiceCore
}

func NewUsersService(core *servicecore.ServiceCore) *ServiceUsers {
	return &ServiceUsers{
		ServiceCore: core,
	}
}

func NewTestUsersService(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *ServiceUsers {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	tokenMaker, _ := token.NewPasetoMaker(config.TokenSymmetricKey)

	serverCore := servicecore.NewServiceCore(config, store, taskDistributor, tokenMaker)

	service := NewUsersService(serverCore)

	return service
}
