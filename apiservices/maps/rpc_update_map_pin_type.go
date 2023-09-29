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

func (server *ServiceMaps) UpdateMapPinType(ctx context.Context, request *pb.UpdateMapPinTypeRequest) (*pb.UpdateMapPinTypeResponse, error) {
	violations := validateUpdateMapPinType(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	err := server.CheckMapAccess(ctx, request.GetMapId(), false)
	if err != nil {
		return nil, err
	}

	argPinType := db.UpdateMapPinTypeParams{
		ID: request.GetPinTypeId(),
		Shape: db.NullPinShape{
			PinShape: converters.ConvertPinShapeToDB(request.GetShape()),
			Valid:    request.Shape != nil,
		},
		BackgroundColor: sql.NullString{
			String: request.GetBackgroundColor(),
			Valid:  request.BackgroundColor != nil,
		},
		BorderColor: sql.NullString{
			String: request.GetBorderColor(),
			Valid:  request.BorderColor != nil,
		},
		IconColor: sql.NullString{
			String: request.GetIconColor(),
			Valid:  request.IconColor != nil,
		},
		Icon: sql.NullString{
			String: request.GetIcon(),
			Valid:  request.Icon != nil,
		},
		IconSize: sql.NullInt32{
			Int32: request.GetIconSize(),
			Valid: request.IconSize != nil,
		},
		Width: sql.NullInt32{
			Int32: request.GetWidth(),
			Valid: request.Width != nil,
		},
	}

	updatedPinType, err := server.Store.UpdateMapPinType(ctx, argPinType)
	if err != nil {
		return nil, err
	}

	rsp := &pb.UpdateMapPinTypeResponse{
		PinType: converters.ConvertMapPinType(updatedPinType),
	}

	return rsp, nil
}

func validateUpdateMapPinType(req *pb.UpdateMapPinTypeRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMapId(req.GetMapId()); err != nil {
		violations = append(violations, e.FieldViolation("map_id", err))
	}

	if err := validator.ValidateUniversalId(req.GetPinTypeId()); err != nil {
		violations = append(violations, e.FieldViolation("pin_type_id", err))
	}

	if req.BackgroundColor != nil {
		if err := validator.ValidateUniversalColor(req.GetBackgroundColor()); err != nil {
			violations = append(violations, e.FieldViolation("background_color", err))
		}
	}

	if req.BorderColor != nil {
		if err := validator.ValidateUniversalColor(req.GetBorderColor()); err != nil {
			violations = append(violations, e.FieldViolation("border_color", err))
		}
	}

	if req.IconColor != nil {
		if err := validator.ValidateUniversalColor(req.GetIconColor()); err != nil {
			violations = append(violations, e.FieldViolation("icon_color", err))
		}
	}

	if req.Icon != nil {
		if err := validator.ValidateUniversalName(req.GetIcon()); err != nil {
			violations = append(violations, e.FieldViolation("icon", err))
		}
	}

	if req.IconSize != nil {
		if err := validator.ValidateUniversalDimension(req.GetIconSize()); err != nil {
			violations = append(violations, e.FieldViolation("icon_size", err))
		}
	}

	if req.Width != nil {
		if err := validator.ValidateUniversalDimension(req.GetWidth()); err != nil {
			violations = append(violations, e.FieldViolation("width", err))
		}
	}

	return violations
}
