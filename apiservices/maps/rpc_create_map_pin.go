package maps

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/api/e"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceMaps) CreateMapPin(ctx context.Context, request *pb.CreateMapPinRequest) (*pb.ViewMapPin, error) {
	violations := validateCreateMapPin(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	_, err := server.CheckEntityTypePermissions(ctx, db.EntityTypeMap, request.GetMapId(), nil)
	if err != nil {
		return nil, err
	}

	argPin := db.CreateMapPinParams{
		MapID: request.GetMapId(),
		Name:  request.GetName(),
		MapPinTypeID: sql.NullInt32{
			Int32: request.GetMapPinTypeId(),
			Valid: request.MapPinTypeId != nil,
		},
		LocationID: sql.NullInt32{
			Int32: request.GetLocationId(),
			Valid: request.LocationId != nil,
		},
		MapLayerID: sql.NullInt32{
			Int32: request.GetMapLayerId(),
			Valid: request.MapLayerId != nil,
		},
		X: request.GetX(),
		Y: request.GetY(),
	}

	newPin, err := server.Store.CreateMapPin(ctx, argPin)
	if err != nil {
		return nil, err
	}

	viewMapPin, err := server.Store.GetMapPinByID(ctx, newPin.ID)
	if err != nil {
		return nil, err
	}

	rsp := converters.ConvertViewMapPin(viewMapPin)

	return rsp, nil
}

func validateCreateMapPin(req *pb.CreateMapPinRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMapId(req.GetMapId()); err != nil {
		violations = append(violations, e.FieldViolation("map_id", err))
	}

	if err := validator.ValidateUniversalName(req.GetName()); err != nil {
		violations = append(violations, e.FieldViolation("name", err))
	}

	if err := validator.ValidateUniversalDimension(req.GetX()); err != nil {
		violations = append(violations, e.FieldViolation("x", err))
	}

	if err := validator.ValidateUniversalDimension(req.GetY()); err != nil {
		violations = append(violations, e.FieldViolation("y", err))
	}

	if req.MapPinTypeId != nil {
		if err := validator.ValidateUniversalId(req.GetMapPinTypeId()); err != nil {
			violations = append(violations, e.FieldViolation("map_pin_type_id", err))
		}
	}

	if req.LocationId != nil {
		if err := validator.ValidateUniversalId(req.GetLocationId()); err != nil {
			violations = append(violations, e.FieldViolation("location_id", err))
		}
	}

	if req.MapLayerId != nil {
		if err := validator.ValidateUniversalId(req.GetMapLayerId()); err != nil {
			violations = append(violations, e.FieldViolation("map_layer_id", err))
		}
	}

	return violations
}
