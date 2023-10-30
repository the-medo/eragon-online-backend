package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertViewImage(viewImage db.ViewImage) *pb.ViewImage {
	pbImage := &pb.ViewImage{
		Id:        viewImage.ID,
		Url:       viewImage.Url,
		BaseUrl:   viewImage.BaseUrl,
		CreatedAt: timestamppb.New(viewImage.CreatedAt),
		UserId:    viewImage.UserID,
		Width:     viewImage.Width,
		Height:    viewImage.Height,
		Tags:      viewImage.Tags,
	}

	if viewImage.ImgGuid.Valid == true {
		pbImage.ImgGuid = viewImage.ImgGuid.UUID.String()
	}

	if viewImage.ImageTypeID.Valid == true {
		pbImage.ImageTypeId = viewImage.ImageTypeID.Int32
	}

	if viewImage.Name.Valid == true {
		pbImage.Name = viewImage.Name.String
	}

	if viewImage.EntityID.Valid == true {
		pbImage.EntityId = &viewImage.EntityID.Int32
	}

	if viewImage.ModuleID.Valid == true {
		pbImage.ModuleId = &viewImage.ModuleID.Int32
	}

	if viewImage.ModuleTypeID.Valid == true {
		pbImage.ModuleId = &viewImage.ModuleTypeID.Int32
	}

	if viewImage.ModuleType.Valid == true {
		convertedModuleType := ConvertModuleTypeToPB(viewImage.ModuleType.ModuleType)
		pbImage.ModuleType = &convertedModuleType
	}

	return pbImage
}
