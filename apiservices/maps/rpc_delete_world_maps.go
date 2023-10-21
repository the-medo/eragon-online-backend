package maps

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

func (server *ServiceMaps) DeleteWorldMap(ctx context.Context, request *pb.DeleteWorldMapRequest) (*emptypb.Empty, error) {
	violations := validateDeleteWorldMap(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	_, err := server.CheckWorldAdmin(ctx, request.GetWorldId(), false)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to delete world map: %v", err)
	}

	err = server.Store.DeleteMapPinsForMap(ctx, request.GetMapId())
	if err != nil {
		return nil, err
	}
	err = server.Store.DeleteMapLayersForMap(ctx, request.GetMapId())
	if err != nil {
		return nil, err
	}

	arg := db.DeleteWorldMapParams{
		WorldID: request.GetWorldId(),
		MapID:   request.GetMapId(),
	}

	err = server.Store.DeleteWorldMap(ctx, arg)
	if err != nil {
		return nil, err
	}

	err = server.Store.DeleteMap(ctx, request.GetMapId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func validateDeleteWorldMap(req *pb.DeleteWorldMapRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateWorldId(req.GetWorldId()); err != nil {
		violations = append(violations, e.FieldViolation("world_id", err))
	}

	if err := validator.ValidateMapId(req.GetMapId()); err != nil {
		violations = append(violations, e.FieldViolation("map_id", err))
	}

	return violations
}
