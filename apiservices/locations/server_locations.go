package locations

import (
	"github.com/the-medo/talebound-backend/apiservices/srv"
	"github.com/the-medo/talebound-backend/pb"
)

type ServerLocations struct {
	pb.UnimplementedLocationsServer
	srv.ServerCore
}

func NewLocationsServer(core srv.ServerCore) ServerLocations {
	return ServerLocations{
		ServerCore: core,
	}
}
