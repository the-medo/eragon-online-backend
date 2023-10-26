package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//TODO: maybe, once generics will help... -_-

func ConvertViewPostToDataPost(viewPost db.ViewPost) *pb.DataPost {
	pbPost := &pb.DataPost{
		Id:         viewPost.ID,
		PostTypeId: viewPost.PostTypeID,
		UserId:     viewPost.UserID,
		Title:      viewPost.Title,
		Content:    viewPost.Content,
		CreatedAt:  timestamppb.New(viewPost.CreatedAt),
		IsDraft:    viewPost.IsDraft,
		IsPrivate:  viewPost.IsPrivate,
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

	if viewPost.ThumbnailImgUrl.Valid == true {
		pbPost.ImageThumbnailUrl = viewPost.ThumbnailImgUrl.String
	}

	return pbPost
}

func ConvertViewPostByPlacementToPost(viewPost db.GetPostsByPlacementRow) *pb.Post {
	pbPost := &pb.Post{
		Post: &pb.DataPost{
			Id:         viewPost.ID,
			PostTypeId: viewPost.PostTypeID,
			UserId:     viewPost.UserID,
			Title:      viewPost.Title,
			Content:    viewPost.Content,
			CreatedAt:  timestamppb.New(viewPost.CreatedAt),
			IsDraft:    viewPost.IsDraft,
			IsPrivate:  viewPost.IsPrivate,
		},
		PostType: &pb.DataPostType{
			Id:         viewPost.PostTypeID,
			Name:       viewPost.PostTypeName,
			Draftable:  viewPost.PostTypeDraftable,
			Privatable: viewPost.PostTypePrivatable,
		},
	}

	if viewPost.DeletedAt.Valid == true {
		pbPost.Post.DeletedAt = timestamppb.New(viewPost.DeletedAt.Time)
	}

	if viewPost.LastUpdatedAt.Valid == true {
		pbPost.Post.LastUpdatedAt = timestamppb.New(viewPost.LastUpdatedAt.Time)
	}

	if viewPost.LastUpdatedUserID.Valid == true {
		pbPost.Post.LastUpdatedUserId = viewPost.LastUpdatedUserID.Int32
	}

	if viewPost.Description.Valid == true {
		pbPost.Post.Description = viewPost.Description.String
	}

	if viewPost.ThumbnailImgID.Valid == true {
		pbPost.Post.ImageThumbnailId = viewPost.ThumbnailImgID.Int32
	}

	if viewPost.ThumbnailImgUrl.Valid == true {
		pbPost.Post.ImageThumbnailUrl = viewPost.ThumbnailImgUrl.String
	}

	return pbPost
}
