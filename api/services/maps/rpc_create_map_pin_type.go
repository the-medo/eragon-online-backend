package maps

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/api/servicecore"
	"github.com/the-medo/talebound-backend/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceMaps) CreateMapPinType(ctx context.Context, request *pb.CreateMapPinTypeRequest) (*pb.MapPinType, error) {
	violations := validateCreateMapPinType(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	_, err := server.CheckModuleIdPermissions(ctx, request.GetModuleId(), &servicecore.ModulePermission{
		NeedsEntityPermission: &[]db.EntityType{db.EntityTypeMap},
	})
	if err != nil {
		return nil, err
	}

	argPinType := db.CreateMapPinTypeParams{
		MapPinTypeGroupID: request.GetMapPinTypeGroupId(),
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
		IsDefault: sql.NullBool{
			Bool:  false,
			Valid: true,
		},
	}

	newPinType, err := server.Store.CreateMapPinType(ctx, argPinType)
	if err != nil {
		return nil, err
	}

	rsp := converters.ConvertMapPinType(newPinType)

	return rsp, nil
}

func validateCreateMapPinType(req *pb.CreateMapPinTypeRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUniversalId(req.GetModuleId()); err != nil {
		violations = append(violations, e.FieldViolation("module_id", err))
	}

	if err := validator.ValidateUniversalId(req.GetMapPinTypeGroupId()); err != nil {
		violations = append(violations, e.FieldViolation("map_pin_type_group_id", err))
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
