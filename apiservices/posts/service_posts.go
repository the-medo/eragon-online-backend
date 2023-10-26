package posts

import (
	"github.com/the-medo/talebound-backend/apiservices/srv"
	"github.com/the-medo/talebound-backend/pb"
)

type ServicePosts struct {
	pb.UnimplementedPostsServer
	*srv.ServiceCore
}

func NewPostsService(core *srv.ServiceCore) *ServicePosts {
	return &ServicePosts{
		ServiceCore: core,
	}
}
