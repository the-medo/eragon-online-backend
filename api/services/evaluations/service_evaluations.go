package evaluations

import (
	"github.com/the-medo/talebound-backend/api/servicecore"
	"github.com/the-medo/talebound-backend/pb"
)

type ServiceEvaluations struct {
	pb.UnimplementedEvaluationsServer
	*servicecore.ServiceCore
}

func NewEvaluationsService(core *servicecore.ServiceCore) *ServiceEvaluations {
	return &ServiceEvaluations{
		ServiceCore: core,
	}
}
