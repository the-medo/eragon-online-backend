package menus

import (
	"github.com/the-medo/talebound-backend/apiservices/srv"
	"github.com/the-medo/talebound-backend/pb"
)

type ServiceMenus struct {
	pb.UnimplementedMenusServer
	*srv.ServiceCore
}

func NewMenusService(core *srv.ServiceCore) *ServiceMenus {
	return &ServiceMenus{
		ServiceCore: core,
	}
}
