package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//TODO: maybe, once generics will help... -_-
// complete copy of ConvertViewPost, but with different name

func ConvertGetPostsRow(viewPost db.GetPostsRow) *pb.Post {
	pbPost := &pb.Post{
		Id:        viewPost.ID,
		UserId:    viewPost.UserID,
		Title:     viewPost.Title,
		Content:   viewPost.Content,
		CreatedAt: timestamppb.New(viewPost.CreatedAt),
		IsDraft:   viewPost.IsDraft,
		IsPrivate: viewPost.IsPrivate,
	}

	if viewPost.DeletedAt.Valid == true {
		pbPost.DeletedAt = timestamppb.New(viewPost.DeletedAt.Time)
	}

	if viewPost.LastUpdatedAt.Valid == true {
		pbPost.LastUpdatedAt = timestamppb.New(viewPost.LastUpdatedAt.Time)
	}

	if viewPost.LastUpdatedUserID.Valid == true {
		pbPost.LastUpdatedUserId = viewPost.LastUpdatedUserID.Int32
	}

	if viewPost.Description.Valid == true {
		pbPost.Description = viewPost.Description.String
	}

	if viewPost.ThumbnailImgID.Valid == true {
		pbPost.ImageThumbnailId = viewPost.ThumbnailImgID.Int32
	}

	return pbPost
}
