package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ConvertViewMap converts a db.ViewMap to pb.ViewMap
func ConvertViewMap(viewMap db.ViewMap) *pb.ViewMap {
	pbMap := &pb.ViewMap{
		Id:        viewMap.ID,
		Title:     viewMap.Title,
		Width:     viewMap.Width,
		Height:    viewMap.Height,
		Tags:      viewMap.Tags,
		CreatedAt: timestamppb.New(viewMap.CreatedAt),
		IsPrivate: viewMap.IsPrivate,
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

	if viewMap.LastUpdatedAt.Valid {
		pbMap.LastUpdatedAt = timestamppb.New(viewMap.LastUpdatedAt.Time)
	}

	if viewMap.LastUpdatedUserID.Valid {
		pbMap.LastUpdatedUserId = viewMap.LastUpdatedUserID.Int32
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
