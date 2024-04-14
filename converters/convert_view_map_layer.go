package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

// ConvertViewMapLayer converts a db.ViewMapLayer to pb.ViewMapLayer
func ConvertViewMapLayer(viewMapLayer db.ViewMapLayer) *pb.ViewMapLayer {
	pbLayer := &pb.ViewMapLayer{
		Id:       viewMapLayer.ID,
		Name:     viewMapLayer.Name,
		MapId:    viewMapLayer.MapID,
		ImageId:  viewMapLayer.ImageID,
		Enabled:  viewMapLayer.Enabled,
		Position: viewMapLayer.Position,
	}

	if viewMapLayer.ImageUrl.Valid {
		pbLayer.ImageUrl = viewMapLayer.ImageUrl.String
	}

	return pbLayer
}
