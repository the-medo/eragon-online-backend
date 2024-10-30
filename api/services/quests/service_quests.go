package quests

import (
	"github.com/the-medo/talebound-backend/api/servicecore"
	"github.com/the-medo/talebound-backend/pb"
)

type ServiceQuests struct {
	pb.UnimplementedQuestsServer
	*servicecore.ServiceCore
}

func NewQuestsService(core *servicecore.ServiceCore) *ServiceQuests {
	return &ServiceQuests{
		ServiceCore: core,
	}
}
