package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertHistoryPost(postHistory db.GetPostHistoryByIdRow, postType db.PostType) *pb.HistoryPost {
	pbHistoryPost := &pb.HistoryPost{
		Post: &pb.DataHistoryPost{
			Id:         postHistory.PostHistoryID,
			PostId:     postHistory.PostID,
			PostTypeId: postHistory.PostTypeID,
			UserId:     postHistory.UserID,
			Title:      postHistory.Title,
			Content:    postHistory.Content,
			CreatedAt:  timestamppb.New(postHistory.CreatedAt),
			IsDraft:    postHistory.IsDraft,
			IsPrivate:  postHistory.IsPrivate,
		},
		PostType: &pb.DataPostType{
			Id:         postType.ID,
			Name:       postType.Name,
			Draftable:  postType.Draftable,
			Privatable: postType.Privatable,
		},
	}

	if postHistory.DeletedAt.Valid == true {
		pbHistoryPost.Post.DeletedAt = timestamppb.New(postHistory.DeletedAt.Time)
	}

	if postHistory.LastUpdatedAt.Valid == true {
		pbHistoryPost.Post.LastUpdatedAt = timestamppb.New(postHistory.LastUpdatedAt.Time)
	}

	if postHistory.LastUpdatedUserID.Valid == true {
		pbHistoryPost.Post.LastUpdatedUserId = postHistory.LastUpdatedUserID.Int32
	}

	if postHistory.Description.Valid == true {
		pbHistoryPost.Post.Description = postHistory.Description.String
	}

	if postHistory.ThumbnailImgID.Valid == true {
		pbHistoryPost.Post.ImageThumbnailId = postHistory.ThumbnailImgID.Int32
	}

	return pbHistoryPost
}
