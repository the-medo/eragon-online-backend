package locations

import (
	"context"
	"github.com/the-medo/talebound-backend/api/e"
	"github.com/the-medo/talebound-backend/apiservices/srv"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServiceLocations struct {
	pb.UnimplementedLocationsServer
	*srv.ServiceCore
}

func (server *ServiceLocations) CheckLocationAccess(ctx context.Context, locationId int32, needsSuperAdmin bool) (*token.Payload, *pb.LocationPlacement, error) {
	var authPayload *token.Payload = nil
	var locationPlacement *pb.LocationPlacement = &pb.LocationPlacement{}

	assignments, err := server.Store.GetLocationAssignments(ctx, locationId)
	if assignments.WorldID > 0 {
		locationPlacement.WorldId = &assignments.WorldID
		authPayload, err = server.CheckWorldAdmin(ctx, assignments.WorldID, needsSuperAdmin)
		if err != nil {
			return nil, nil, status.Errorf(codes.PermissionDenied, "failed to get location access - not world admin: %v", err)
		}
	}

	_, err = server.AuthorizeUserCookie(ctx)
	if err != nil {
		return nil, nil, e.UnauthenticatedError(err)
	}

	return authPayload, locationPlacement, nil
}

func NewLocationsService(core *srv.ServiceCore) *ServiceLocations {
	return &ServiceLocations{
		ServiceCore: core,
	}
}
