package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

func ConvertModuleTypeTag(tag db.ModuleTypeTagsAvailable) *pb.Tag {
	pbTag := &pb.Tag{
		Id:  tag.ID,
		Tag: tag.Tag,
	}

	return pbTag
}

func ConvertViewTag(tag db.ViewModuleTypeTagsAvailable) *pb.ViewTag {
	pbTag := &pb.ViewTag{
		Id:         tag.ID,
		Tag:        tag.Tag,
		ModuleType: ConvertModuleTypeToPB(tag.ModuleType),
		Count:      tag.Count,
	}
	return pbTag
}

func ConvertViewTagToTag(tag db.ViewModuleTypeTagsAvailable) *pb.Tag {
	pbTag := &pb.Tag{
		Id:  tag.ID,
		Tag: tag.Tag,
	}
	return pbTag
}
