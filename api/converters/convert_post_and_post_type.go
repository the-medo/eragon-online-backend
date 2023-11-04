package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

func ConvertPostAndPostType(post db.ViewPost, postType db.PostType) *pb.Post {
	pbPost := &pb.Post{
		Post: ConvertViewPostToDataPost(post),
		PostType: &pb.DataPostType{
			Id:         postType.ID,
			Name:       postType.Name,
			Draftable:  postType.Draftable,
			Privatable: postType.Privatable,
		},
	}

	return pbPost
}
