package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

func ConvertPostType(postType db.PostType) *pb.DataPostType {
	pbPostType := &pb.DataPostType{
		Id:         postType.ID,
		Name:       postType.Name,
		Draftable:  postType.Draftable,
		Privatable: postType.Privatable,
	}

	return pbPostType
}
