package entities

import (
	"github.com/the-medo/talebound-backend/apiservices/srv"
	"github.com/the-medo/talebound-backend/pb"
)

type ServiceEntities struct {
	pb.UnimplementedEntitiesServer
	*srv.ServiceCore
}

func NewEntitiesService(core *srv.ServiceCore) *ServiceEntities {
	return &ServiceEntities{
		ServiceCore: core,
	}
}
