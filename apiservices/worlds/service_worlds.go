package worlds

import (
	"github.com/the-medo/talebound-backend/apiservices/servicecore"
	"github.com/the-medo/talebound-backend/pb"
)

type ServiceWorlds struct {
	pb.UnimplementedWorldsServer
	*servicecore.ServiceCore
}

func NewWorldsService(core *servicecore.ServiceCore) *ServiceWorlds {
	return &ServiceWorlds{
		ServiceCore: core,
	}
}
