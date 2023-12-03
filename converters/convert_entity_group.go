package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

func ConvertEntityGroup(entityGroup db.EntityGroup) *pb.EntityGroup {
	pbGroup := &pb.EntityGroup{
		Id: entityGroup.ID,
	}

	if entityGroup.Name.Valid {
		pbGroup.Name = &entityGroup.Name.String
	}
	if entityGroup.Description.Valid {
		pbGroup.Description = &entityGroup.Description.String
	}

	return pbGroup
}
