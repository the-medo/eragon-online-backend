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

func ConvertViewTag(tag db.ViewWorldTagsAvailable) *pb.ViewTag {
	pbTag := &pb.ViewTag{
		Id:    tag.ID,
		Tag:   tag.Tag,
		Count: tag.Count,
	}
	return pbTag
}

func ConvertViewTagToTag(tag db.ViewWorldTagsAvailable) *pb.Tag {
	pbTag := &pb.Tag{
		Id:  tag.ID,
		Tag: tag.Tag,
	}
	return pbTag
}
