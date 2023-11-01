package maps

import (
	"context"
	"github.com/the-medo/talebound-backend/api/e"
	"github.com/the-medo/talebound-backend/apiservices/srv"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServiceMaps struct {
	pb.UnimplementedMapsServer
	*srv.ServiceCore
}

func (server *ServiceMaps) CheckMapAccess(ctx context.Context, mapId int32, needsSuperAdmin bool) error {
	assignments, err := server.Store.GetMapAssignments(ctx, mapId)
	if assignments.WorldID > 0 {
		_, err = server.CheckWorldPermissions(ctx, assignments.WorldID, needsSuperAdmin)
		if err != nil {
			return status.Errorf(codes.PermissionDenied, "failed to get map access - not world admin: %v", err)
		}
	}

	_, err = server.AuthorizeUserCookie(ctx)
	if err != nil {
		return e.UnauthenticatedError(err)
	}

	return nil
}

func NewMapsService(core *srv.ServiceCore) *ServiceMaps {
	return &ServiceMaps{
		ServiceCore: core,
	}
}
