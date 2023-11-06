package maps

import (
	"github.com/the-medo/talebound-backend/apiservices/servicecore"
	"github.com/the-medo/talebound-backend/pb"
)

type ServiceMaps struct {
	pb.UnimplementedMapsServer
	*servicecore.ServiceCore
}

func NewMapsService(core *servicecore.ServiceCore) *ServiceMaps {
	return &ServiceMaps{
		ServiceCore: core,
	}
}
