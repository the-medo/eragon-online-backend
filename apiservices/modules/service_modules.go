package modules

import (
	"github.com/the-medo/talebound-backend/apiservices/srv"
	"github.com/the-medo/talebound-backend/pb"
)

type ServiceModules struct {
	pb.UnimplementedModulesServer
	*srv.ServiceCore
}

func NewModulesService(core *srv.ServiceCore) *ServiceModules {
	return &ServiceModules{
		ServiceCore: core,
	}
}
