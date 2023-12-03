package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

// ConvertViewMap converts a db.ViewMap to pb.ViewMap
func ConvertViewMap(viewMap db.ViewMap) *pb.ViewMap {
	pbMap := &pb.ViewMap{
		Id:     viewMap.ID,
		Name:   viewMap.Name,
		Width:  viewMap.Width,
		Height: viewMap.Height,
		Tags:   viewMap.Tags,
	}

	if viewMap.Type.Valid {
		pbMap.Type = &viewMap.Type.String
	}

	if viewMap.Description.Valid {
		pbMap.Description = &viewMap.Description.String
	}

	if viewMap.ThumbnailImageID.Valid {
		pbMap.ThumbnailImageId = &viewMap.ThumbnailImageID.Int32
	}

	if viewMap.ThumbnailImageUrl.Valid {
		pbMap.ThumbnailImageUrl = &viewMap.ThumbnailImageUrl.String
	}

	if viewMap.EntityID.Valid == true {
		pbMap.EntityId = viewMap.EntityID.Int32
	}

	if viewMap.ModuleID.Valid == true {
		pbMap.ModuleId = viewMap.ModuleID.Int32
	}

	if viewMap.ModuleTypeID.Valid == true {
		pbMap.ModuleId = viewMap.ModuleTypeID.Int32
	}

	if viewMap.ModuleType.Valid == true {
		convertedModuleType := ConvertModuleTypeToPB(viewMap.ModuleType.ModuleType)
		pbMap.ModuleType = convertedModuleType
	}

	return pbMap
}
