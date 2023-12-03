package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//TODO: maybe, once generics will help... -_-
// complete copy of ConvertViewPost, but with different name

func ConvertGetImagesRow(viewImage db.GetImagesRow) *pb.Image {
	pbImage := &pb.Image{
		Id:        viewImage.ID,
		UserId:    viewImage.UserID,
		CreatedAt: timestamppb.New(viewImage.CreatedAt),
		Width:     viewImage.Width,
		Height:    viewImage.Height,
		Url:       viewImage.Url,
		BaseUrl:   viewImage.BaseUrl,
	}

	if viewImage.Name.Valid == true {
		pbImage.Name = viewImage.Name.String
	}

	if viewImage.ImgGuid.Valid == true {
		pbImage.ImgGuid = viewImage.ImgGuid.UUID.String()
	}

	if viewImage.ImageTypeID.Valid == true {
		pbImage.ImageTypeId = viewImage.ImageTypeID.Int32
	}

	return pbImage
}
