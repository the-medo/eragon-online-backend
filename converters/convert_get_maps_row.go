package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

func ConvertGetMapsRow(viewMap db.GetMapsRow) *pb.Map {
	pbMap := &pb.Map{
		Id:     viewMap.ID,
		Name:   viewMap.Name,
		Width:  viewMap.Width,
		Height: viewMap.Height,
	}

	if viewMap.Description.Valid {
		pbMap.Description = &viewMap.Description.String
	}

	if viewMap.Type.Valid {
		pbMap.Type = &viewMap.Type.String
	}

	if viewMap.ThumbnailImageID.Valid {
		pbMap.ThumbnailImageId = &viewMap.ThumbnailImageID.Int32
	}

	return pbMap
}
