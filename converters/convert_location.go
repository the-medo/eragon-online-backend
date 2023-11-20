package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

func ConvertLocation(viewLocation db.Location) *pb.Location {
	pbLocation := &pb.Location{
		Id:   viewLocation.ID,
		Name: viewLocation.Name,
	}

	if viewLocation.Description.Valid == true {
		pbLocation.Description = &viewLocation.Description.String
	}

	if viewLocation.PostID.Valid == true {
		pbLocation.PostId = &viewLocation.PostID.Int32
	}

	if viewLocation.ThumbnailImageID.Valid == true {
		pbLocation.ThumbnailImageId = &viewLocation.ThumbnailImageID.Int32
	}

	return pbLocation
}
