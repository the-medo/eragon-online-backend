package entities

import (
	"github.com/the-medo/talebound-backend/apiservices/servicecore"
	"github.com/the-medo/talebound-backend/pb"
)

type ServiceEntities struct {
	pb.UnimplementedEntitiesServer
	*servicecore.ServiceCore
}

func NewEntitiesService(core *servicecore.ServiceCore) *ServiceEntities {
	return &ServiceEntities{
		ServiceCore: core,
	}
}
