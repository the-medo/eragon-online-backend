package menus

import (
	"github.com/the-medo/talebound-backend/apiservices/servicecore"
	"github.com/the-medo/talebound-backend/pb"
)

type ServiceMenus struct {
	pb.UnimplementedMenusServer
	*servicecore.ServiceCore
}

func NewMenusService(core *servicecore.ServiceCore) *ServiceMenus {
	return &ServiceMenus{
		ServiceCore: core,
	}
}
