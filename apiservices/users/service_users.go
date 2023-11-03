package users

import (
	"github.com/the-medo/talebound-backend/apiservices/srv"
	"github.com/the-medo/talebound-backend/pb"
)

type ServiceUsers struct {
	pb.UnimplementedUsersServer
	*srv.ServiceCore
}

func NewUsersService(core *srv.ServiceCore) *ServiceUsers {
	return &ServiceUsers{
		ServiceCore: core,
	}
}
