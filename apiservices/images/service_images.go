package images

import (
	"github.com/the-medo/talebound-backend/apiservices/srv"
	"github.com/the-medo/talebound-backend/pb"
)

type ServiceImages struct {
	pb.UnimplementedImagesServer
	*srv.ServiceCore
}

func NewImagesService(core *srv.ServiceCore) *ServiceImages {
	return &ServiceImages{
		ServiceCore: core,
	}
}
