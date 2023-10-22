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

func (server *ServiceMaps) CreateMapPinType(ctx context.Context, request *pb.CreateMapPinTypeRequest) (*pb.MapPinType, error) {
	violations := validateCreateMapPinType(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	err := server.CheckMapAccess(ctx, request.GetMapId(), false)
	if err != nil {
		return nil, err
	}

	mapPinTypeGroupId, err := server.Store.GetMapPinTypeGroupIdForMap(ctx, request.GetMapId())
	if err != nil {
		return nil, err
	}

	argPinType := db.CreateMapPinTypeParams{
		MapPinTypeGroupID: mapPinTypeGroupId,
		Shape:             converters.ConvertPinShapeToDB(request.GetShape()),
		BackgroundColor: sql.NullString{
			String: request.GetBackgroundColor(),
			Valid:  true,
		},
		BorderColor: sql.NullString{
			String: request.GetBorderColor(),
			Valid:  true,
		},
		IconColor: sql.NullString{
			String: request.GetIconColor(),
			Valid:  true,
		},
		Icon: sql.NullString{
			String: request.GetIcon(),
			Valid:  true,
		},
		IconSize: sql.NullInt32{
			Int32: request.GetIconSize(),
			Valid: true,
		},
		Width: sql.NullInt32{
			Int32: request.GetWidth(),
			Valid: true,
		},
		Section: request.GetSection(),
	}

	newPinType, err := server.Store.CreateMapPinType(ctx, argPinType)
	if err != nil {
		return nil, err
	}

	rsp := converters.ConvertMapPinType(newPinType)

	return rsp, nil
}

func validateCreateMapPinType(req *pb.CreateMapPinTypeRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMapId(req.GetMapId()); err != nil {
		violations = append(violations, e.FieldViolation("map_id", err))
	}

	if err := validator.ValidateUniversalColor(req.GetBackgroundColor()); err != nil {
		violations = append(violations, e.FieldViolation("background_color", err))
	}

	if err := validator.ValidateUniversalColor(req.GetBorderColor()); err != nil {
		violations = append(violations, e.FieldViolation("border_color", err))
	}

	if err := validator.ValidateUniversalColor(req.GetIconColor()); err != nil {
		violations = append(violations, e.FieldViolation("icon_color", err))
	}

	if err := validator.ValidateUniversalName(req.GetIcon()); err != nil {
		violations = append(violations, e.FieldViolation("icon", err))
	}

	if err := validator.ValidateUniversalDimension(req.GetIconSize()); err != nil {
		violations = append(violations, e.FieldViolation("icon_size", err))
	}

	if err := validator.ValidateUniversalDimension(req.GetWidth()); err != nil {
		violations = append(violations, e.FieldViolation("width", err))
	}

	return violations
}
