package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

func ConvertEntityGroupContent(content db.EntityGroupContent) *pb.EntityGroupContent {
	pbContent := &pb.EntityGroupContent{
		Id:            content.ID,
		EntityGroupId: content.EntityGroupID,
		Position:      content.Position,
	}

	if content.ContentEntityID.Valid {
		pbContent.ContentEntityId = &content.ContentEntityID.Int32
	}
	if content.ContentEntityGroupID.Valid {
		pbContent.ContentEntityGroupId = &content.ContentEntityGroupID.Int32
	}

	return pbContent
}
