package auth

import (
	"github.com/the-medo/talebound-backend/apiservices/srv"
	"github.com/the-medo/talebound-backend/pb"
)

type ServiceAuth struct {
	pb.UnimplementedAuthServer
	*srv.ServiceCore
}

func NewAuthService(core *srv.ServiceCore) *ServiceAuth {
	return &ServiceAuth{
		ServiceCore: core,
	}
}
