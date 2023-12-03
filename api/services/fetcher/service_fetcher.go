package fetcher

import (
	"github.com/the-medo/talebound-backend/api/servicecore"
	"github.com/the-medo/talebound-backend/pb"
)

type ServiceFetcher struct {
	pb.UnimplementedFetcherServer
	*servicecore.ServiceCore
}

func NewFetcherService(core *servicecore.ServiceCore) *ServiceFetcher {
	return &ServiceFetcher{
		ServiceCore: core,
	}
}
