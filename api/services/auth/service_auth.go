package auth

import (
	"github.com/the-medo/talebound-backend/api/servicecore"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/token"
	"github.com/the-medo/talebound-backend/util"
	"github.com/the-medo/talebound-backend/worker"
	"testing"
	"time"
)

type ServiceAuth struct {
	pb.UnimplementedAuthServer
	*servicecore.ServiceCore
}

func NewAuthService(core *servicecore.ServiceCore) *ServiceAuth {
	return &ServiceAuth{
		ServiceCore: core,
	}
}

func NewTestAuthService(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *ServiceAuth {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	tokenMaker, _ := token.NewPasetoMaker(config.TokenSymmetricKey)

	serverCore := servicecore.NewServiceCore(config, store, taskDistributor, tokenMaker)

	service := NewAuthService(serverCore)

	return service
}
