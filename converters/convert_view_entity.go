package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

func ConvertViewEntity(dbEntity db.ViewEntity) *pb.ViewEntity {
	pbEntity := &pb.ViewEntity{
		Id:           dbEntity.ID,
		ModuleId:     dbEntity.ModuleID,
		ModuleTypeId: dbEntity.ModuleTypeID,
		Type:         entityTypeToPB[dbEntity.Type],
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

	if dbEntity.ModuleType.Valid == true {
		convertedModuleType := ConvertModuleTypeToPB(dbEntity.ModuleType.ModuleType)
		pbEntity.ModuleType = convertedModuleType
	}

	return pbEntity
}
