package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

// ConvertMapPinType converts a db.MapPinType to pb.MapPinType
func ConvertMapPinType(mapPinType db.MapPinType) *pb.MapPinType {
	pbType := &pb.MapPinType{
		Id:                mapPinType.ID,
		MapPinTypeGroupId: mapPinType.MapPinTypeGroupID,
		Shape:             ConvertPinShapeToPB(mapPinType.Shape),
		IsDefault:         mapPinType.IsDefault,
	}

	if mapPinType.BackgroundColor.Valid {
		pbType.BackgroundColor = mapPinType.BackgroundColor.String
	}

	if mapPinType.BorderColor.Valid {
		pbType.BorderColor = mapPinType.BorderColor.String
	}

	if mapPinType.IconColor.Valid {
		pbType.IconColor = mapPinType.IconColor.String
	}

	if mapPinType.Icon.Valid {
		pbType.Icon = mapPinType.Icon.String
	}

	if mapPinType.IconSize.Valid {
		pbType.IconSize = mapPinType.IconSize.Int32
	}

	if mapPinType.Width.Valid {
		pbType.Width = mapPinType.Width.Int32
	}

	return pbType
}
