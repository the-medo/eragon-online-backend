package entities

import (
	"context"
	"github.com/the-medo/talebound-backend/api/e"
	"github.com/the-medo/talebound-backend/apiservices/srv"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServiceEntities struct {
	pb.UnimplementedEntitiesServer
	*srv.ServiceCore
}

func (server *ServiceEntities) CheckEntityGroupAccess(ctx context.Context, entityGroupId int32, needsSuperAdmin bool) error {
	menuId, err := server.Store.GetMenuIdOfEntityGroup(ctx, entityGroupId)
	_, err = server.CheckMenuAdmin(ctx, menuId, needsSuperAdmin)
	if err != nil {
		return status.Errorf(codes.PermissionDenied, "failed to get entity group access - not menu admin: %v", err)
	}

	_, err = server.AuthorizeUserCookie(ctx)
	if err != nil {
		return e.UnauthenticatedError(err)
	}

	return nil
}

func NewEntitiesService(core *srv.ServiceCore) *ServiceEntities {
	return &ServiceEntities{
		ServiceCore: core,
	}
}
