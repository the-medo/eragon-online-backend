package worlds

import (
	"github.com/the-medo/talebound-backend/apiservices/srv"
	"github.com/the-medo/talebound-backend/pb"
)

type ServiceWorlds struct {
	pb.UnimplementedWorldsServer
	*srv.ServiceCore
}

func NewWorldsService(core *srv.ServiceCore) *ServiceWorlds {
	return &ServiceWorlds{
		ServiceCore: core,
	}
}
