package modules

import (
	"github.com/the-medo/talebound-backend/apiservices/servicecore"
	"github.com/the-medo/talebound-backend/pb"
)

type ServiceModules struct {
	pb.UnimplementedModulesServer
	*servicecore.ServiceCore
}

func NewModulesService(core *servicecore.ServiceCore) *ServiceModules {
	return &ServiceModules{
		ServiceCore: core,
	}
}
