package tags

import (
	"github.com/the-medo/talebound-backend/api/servicecore"
	"github.com/the-medo/talebound-backend/pb"
)

type ServiceTags struct {
	pb.UnimplementedTagsServer
	*servicecore.ServiceCore
}

func NewTagsService(core *servicecore.ServiceCore) *ServiceTags {
	return &ServiceTags{
		ServiceCore: core,
	}
}
