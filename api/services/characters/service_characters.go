package characters

import (
	"github.com/the-medo/talebound-backend/api/servicecore"
	"github.com/the-medo/talebound-backend/pb"
)

type ServiceCharacters struct {
	pb.UnimplementedCharactersServer
	*servicecore.ServiceCore
}

func NewCharactersService(core *servicecore.ServiceCore) *ServiceCharacters {
	return &ServiceCharacters{
		ServiceCore: core,
	}
}
