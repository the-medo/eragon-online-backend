package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

func ConvertGetLocationsRow(viewLocation db.GetLocationsRow) *pb.Location {
	pbLocation := &pb.Location{
		Id:   viewLocation.ID,
		Name: viewLocation.Name,
	}

	if viewLocation.Description.Valid {
		pbLocation.Description = &viewLocation.Description.String
	}

	if viewLocation.PostID.Valid {
		pbLocation.PostId = &viewLocation.PostID.Int32
	}

	if viewLocation.ThumbnailImageID.Valid {
		pbLocation.ThumbnailImageId = &viewLocation.ThumbnailImageID.Int32
	}

	return pbLocation
}
