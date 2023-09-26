package maps

import (
	"github.com/the-medo/talebound-backend/apiservices/srv"
	"github.com/the-medo/talebound-backend/pb"
)

type ServiceMaps struct {
	pb.UnimplementedMapsServer
	*srv.ServiceCore
}

func NewMapsService(core *srv.ServiceCore) *ServiceMaps {
	return &ServiceMaps{
		ServiceCore: core,
	}
}
