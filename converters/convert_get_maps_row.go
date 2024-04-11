package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertGetMapsRow(viewMap db.GetMapsRow) *pb.Map {
	pbMap := &pb.Map{
		Id:        viewMap.ID,
		Title:     viewMap.Title,
		Width:     viewMap.Width,
		Height:    viewMap.Height,
		CreatedAt: timestamppb.New(viewMap.CreatedAt),
		IsPrivate: viewMap.IsPrivate,
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

	if viewMap.LastUpdatedAt.Valid {
		pbMap.LastUpdatedAt = timestamppb.New(viewMap.LastUpdatedAt.Time)
	}

	if viewMap.LastUpdatedUserID.Valid {
		pbMap.LastUpdatedUserId = viewMap.LastUpdatedUserID.Int32
	}

	return pbMap
}
