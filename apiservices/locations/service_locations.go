package locations

import (
	"github.com/the-medo/talebound-backend/apiservices/servicecore"
	"github.com/the-medo/talebound-backend/pb"
)

type ServiceLocations struct {
	pb.UnimplementedLocationsServer
	*servicecore.ServiceCore
}

func NewLocationsService(core *servicecore.ServiceCore) *ServiceLocations {
	return &ServiceLocations{
		ServiceCore: core,
	}
}
