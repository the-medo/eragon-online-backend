package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

func ConvertViewLocation(viewLocation db.ViewLocation) *pb.ViewLocation {
	pbLocation := &pb.ViewLocation{
		Id:   viewLocation.ID,
		Name: viewLocation.Name,
	}

	if viewLocation.Description.Valid == true {
		pbLocation.Description = &viewLocation.Description.String
	}

	if viewLocation.PostID.Valid == true {
		pbLocation.PostId = &viewLocation.PostID.Int32
	}

	if viewLocation.PostTitle.Valid == true {
		pbLocation.PostTitle = &viewLocation.PostTitle.String
	}

	if viewLocation.ThumbnailImageID.Valid == true {
		pbLocation.ThumbnailImageId = &viewLocation.ThumbnailImageID.Int32
	}

	if viewLocation.ThumbnailImageUrl.Valid == true {
		pbLocation.ThumbnailImageUrl = &viewLocation.ThumbnailImageUrl.String
	}

	return pbLocation
}
