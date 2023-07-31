package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

func ConvertTag(tag db.WorldTagsAvailable) *pb.Tag {
	pbTag := &pb.Tag{
		Id:  tag.ID,
		Tag: tag.Tag,
	}

	return pbTag
}
