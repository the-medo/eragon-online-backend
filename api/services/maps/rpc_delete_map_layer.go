package maps

import (
	"context"
	"database/sql"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *ServiceMaps) DeleteMapLayer(ctx context.Context, request *pb.DeleteMapLayerRequest) (*emptypb.Empty, error) {
	violations := validateDeleteMapLayer(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	_, err := server.CheckEntityTypePermissions(ctx, db.EntityTypeMap, request.GetMapId(), nil)
	if err != nil {
		return nil, err
	}

	mapLayer, err := server.Store.GetMapLayerByID(ctx, request.GetLayerId())
	if err != nil {
		return nil, err
	}
	if mapLayer.IsMain {
		return nil, status.Errorf(codes.PermissionDenied, "cannot delete main map layer")
	}

	err = server.Store.DeleteMapPinsForMapLayer(ctx, sql.NullInt32{
		Int32: request.GetLayerId(),
		Valid: true,
	})

	err = server.Store.DeleteMapLayer(ctx, request.GetLayerId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func validateDeleteMapLayer(req *pb.DeleteMapLayerRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMapId(req.GetMapId()); err != nil {
		violations = append(violations, e.FieldViolation("map_id", err))
	}

	if err := validator.ValidateUniversalId(req.GetLayerId()); err != nil {
		violations = append(violations, e.FieldViolation("layer_id", err))
	}

	return violations
}
