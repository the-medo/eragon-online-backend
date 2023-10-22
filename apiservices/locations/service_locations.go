package locations

import (
	"context"
	"github.com/the-medo/talebound-backend/api/e"
	"github.com/the-medo/talebound-backend/apiservices/srv"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServiceLocations struct {
	pb.UnimplementedLocationsServer
	*srv.ServiceCore
}

func (server *ServiceLocations) CheckLocationAccess(ctx context.Context, locationId int32, needsSuperAdmin bool) error {
	assignments, err := server.Store.GetLocationAssignments(ctx, locationId)
	if assignments.WorldID > 0 {
		_, err = server.CheckWorldAdmin(ctx, assignments.WorldID, needsSuperAdmin)
		if err != nil {
			return status.Errorf(codes.PermissionDenied, "failed to get location access - not world admin: %v", err)
		}
	}

	_, err = server.AuthorizeUserCookie(ctx)
	if err != nil {
		return e.UnauthenticatedError(err)
	}

	return nil
}

func NewLocationsService(core *srv.ServiceCore) *ServiceLocations {
	return &ServiceLocations{
		ServiceCore: core,
	}
}
