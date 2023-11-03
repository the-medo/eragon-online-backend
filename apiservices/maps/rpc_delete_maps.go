package maps

import (
	"context"
	"github.com/the-medo/talebound-backend/api/e"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *ServiceMaps) DeleteMap(ctx context.Context, request *pb.DeleteMapRequest) (*emptypb.Empty, error) {
	violations := validateDeleteWorldMap(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	_, err := server.CheckEntityTypePermissions(ctx, db.EntityTypeMap, request.GetMapId(), nil)
	if err != nil {
		return nil, err
	}

	err = server.Store.DeleteMap(ctx, request.GetMapId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func validateDeleteWorldMap(req *pb.DeleteMapRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMapId(req.GetMapId()); err != nil {
		violations = append(violations, e.FieldViolation("map_id", err))
	}

	return violations
}
