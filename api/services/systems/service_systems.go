package systems

import (
	"github.com/the-medo/talebound-backend/api/servicecore"
	"github.com/the-medo/talebound-backend/pb"
)

type ServiceSystems struct {
	pb.UnimplementedSystemsServer
	*servicecore.ServiceCore
}

func NewSystemsService(core *servicecore.ServiceCore) *ServiceSystems {
	return &ServiceSystems{
		ServiceCore: core,
	}
}
