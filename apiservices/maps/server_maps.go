package maps

import (
	"github.com/the-medo/talebound-backend/apiservices/srv"
	"github.com/the-medo/talebound-backend/pb"
)

type ServerMaps struct {
	pb.UnimplementedMapsServer
	*srv.ServerCore
}

func NewMapsServer(core *srv.ServerCore) *ServerMaps {
	return &ServerMaps{
		ServerCore: core,
	}
}
