package posts

import (
	"github.com/the-medo/talebound-backend/apiservices/servicecore"
	"github.com/the-medo/talebound-backend/pb"
)

type ServicePosts struct {
	pb.UnimplementedPostsServer
	*servicecore.ServiceCore
}

func NewPostsService(core *servicecore.ServiceCore) *ServicePosts {
	return &ServicePosts{
		ServiceCore: core,
	}
}
