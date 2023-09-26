package locations

import (
	"github.com/the-medo/talebound-backend/apiservices/srv"
	"github.com/the-medo/talebound-backend/pb"
)

type ServiceLocations struct {
	pb.UnimplementedLocationsServer
	*srv.ServiceCore
}

func NewLocationsService(core *srv.ServiceCore) *ServiceLocations {
	return &ServiceLocations{
		ServiceCore: core,
	}
}
