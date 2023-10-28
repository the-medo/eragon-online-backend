package tags

import (
	"github.com/the-medo/talebound-backend/apiservices/srv"
	"github.com/the-medo/talebound-backend/pb"
)

type ServiceTags struct {
	pb.UnimplementedTagsServer
	*srv.ServiceCore
}

func NewTagsService(core *srv.ServiceCore) *ServiceTags {
	return &ServiceTags{
		ServiceCore: core,
	}
}
