package evaluations

import (
	"github.com/the-medo/talebound-backend/apiservices/srv"
	"github.com/the-medo/talebound-backend/pb"
)

type ServiceEvaluations struct {
	pb.UnimplementedEvaluationsServer
	*srv.ServiceCore
}

func NewEvaluationsService(core *srv.ServiceCore) *ServiceEvaluations {
	return &ServiceEvaluations{
		ServiceCore: core,
	}
}
