package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertPostHistory(postHistory db.GetPostHistoryByIdRow) *pb.PostHistory {
	pbHistoryPost := &pb.PostHistory{
		Id:        postHistory.PostHistoryID,
		PostId:    postHistory.PostID,
		UserId:    postHistory.UserID,
		Title:     postHistory.Title,
		Content:   postHistory.Content,
		CreatedAt: timestamppb.New(postHistory.CreatedAt),
		IsDraft:   postHistory.IsDraft,
		IsPrivate: postHistory.IsPrivate,
	}

	if postHistory.DeletedAt.Valid == true {
		pbHistoryPost.DeletedAt = timestamppb.New(postHistory.DeletedAt.Time)
	}

	if postHistory.LastUpdatedAt.Valid == true {
		pbHistoryPost.LastUpdatedAt = timestamppb.New(postHistory.LastUpdatedAt.Time)
	}

	if postHistory.LastUpdatedUserID.Valid == true {
		pbHistoryPost.LastUpdatedUserId = postHistory.LastUpdatedUserID.Int32
	}

	if postHistory.Description.Valid == true {
		pbHistoryPost.Description = postHistory.Description.String
	}

	if postHistory.ThumbnailImgID.Valid == true {
		pbHistoryPost.ImageThumbnailId = postHistory.ThumbnailImgID.Int32
	}

	return pbHistoryPost
}
