package images

import (
	"github.com/the-medo/talebound-backend/apiservices/servicecore"
	"github.com/the-medo/talebound-backend/pb"
)

type ServiceImages struct {
	pb.UnimplementedImagesServer
	*servicecore.ServiceCore
}

func NewImagesService(core *servicecore.ServiceCore) *ServiceImages {
	return &ServiceImages{
		ServiceCore: core,
	}
}
