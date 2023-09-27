package locations

import (
	"context"
	"github.com/the-medo/talebound-backend/api/e"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *ServiceLocations) DeleteWorldLocation(ctx context.Context, request *pb.DeleteWorldLocationRequest) (*emptypb.Empty, error) {
	violations := validateDeleteWorldLocation(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	_, err := server.CheckWorldAdmin(ctx, request.GetWorldId(), false)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to add world tag: %v", err)
	}

	arg := db.DeleteWorldLocationParams{
		WorldID:    request.GetWorldId(),
		LocationID: request.GetLocationId(),
	}

	err = server.Store.DeleteWorldLocation(ctx, arg)
	if err != nil {
		return nil, err
	}

	err = server.Store.DeleteLocation(ctx, request.GetLocationId())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func validateDeleteWorldLocation(req *pb.DeleteWorldLocationRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateWorldId(req.GetWorldId()); err != nil {
		violations = append(violations, e.FieldViolation("world_id", err))
	}

	if err := validator.ValidateLocationId(req.GetLocationId()); err != nil {
		violations = append(violations, e.FieldViolation("location_id", err))
	}
	return violations
}
