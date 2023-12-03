package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

// ConvertViewMapPin converts a db.ViewMapPin to pb.ViewMapPin
func ConvertViewMapPin(viewMapPin db.ViewMapPin) *pb.ViewMapPin {
	pbPin := &pb.ViewMapPin{
		Id:    viewMapPin.ID,
		Name:  viewMapPin.Name,
		MapId: viewMapPin.MapID,
		X:     viewMapPin.X,
		Y:     viewMapPin.Y,
	}

	if viewMapPin.MapPinTypeID.Valid {
		pbPin.MapPinTypeId = &viewMapPin.MapPinTypeID.Int32
	}

	if viewMapPin.LocationID.Valid {
		pbPin.LocationId = &viewMapPin.LocationID.Int32
	}

	if viewMapPin.MapLayerID.Valid {
		pbPin.MapLayerId = &viewMapPin.MapLayerID.Int32
	}

	if viewMapPin.LocationName.Valid {
		pbPin.LocationName = &viewMapPin.LocationName.String
	}

	if viewMapPin.LocationPostID.Valid {
		pbPin.LocationPostId = &viewMapPin.LocationPostID.Int32
	}

	if viewMapPin.LocationDescription.Valid {
		pbPin.LocationDescription = &viewMapPin.LocationDescription.String
	}

	if viewMapPin.LocationThumbnailImageID.Valid {
		pbPin.LocationThumbnailImageId = &viewMapPin.LocationThumbnailImageID.Int32
	}

	if viewMapPin.LocationThumbnailImageUrl.Valid {
		pbPin.LocationThumbnailImageUrl = &viewMapPin.LocationThumbnailImageUrl.String
	}

	return pbPin
}
