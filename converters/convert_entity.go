package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

func ConvertEntity(dbEntity db.Entity) *pb.Entity {
	pbEntity := &pb.Entity{
		Id:       dbEntity.ID,
		ModuleId: dbEntity.ModuleID,
		Type:     entityTypeToPB[dbEntity.Type],
	}

	if dbEntity.PostID.Valid {
		pbEntity.PostId = &dbEntity.PostID.Int32
	}

	if dbEntity.MapID.Valid {
		pbEntity.MapId = &dbEntity.MapID.Int32
	}

	if dbEntity.LocationID.Valid {
		pbEntity.LocationId = &dbEntity.LocationID.Int32
	}

	if dbEntity.ImageID.Valid {
		pbEntity.ImageId = &dbEntity.ImageID.Int32
	}

	return pbEntity
}
