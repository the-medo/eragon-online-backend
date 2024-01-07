package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

func ConvertEntityGroup(entityGroup db.EntityGroup) *pb.EntityGroup {
	pbGroup := &pb.EntityGroup{
		Id:        entityGroup.ID,
		Style:     pb.EntityGroupStyle_ENTITY_GROUP_STYLE_UNKNOWN,
		Direction: pb.EntityGroupDirection_ENTITY_GROUP_DIRECTION_UNKNOWN,
	}

	if entityGroup.Name.Valid {
		pbGroup.Name = &entityGroup.Name.String
	}

	if entityGroup.Description.Valid {
		pbGroup.Description = &entityGroup.Description.String
	}

	if entityGroup.Style.Valid {
		pbGroup.Style = ConvertEntityGroupStyleToPB(entityGroup.Style.String)
	}

	if entityGroup.Direction.Valid {
		pbGroup.Direction = ConvertEntityGroupDirectionToPB(entityGroup.Direction.String)
	}

	return pbGroup
}
