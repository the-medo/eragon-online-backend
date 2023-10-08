package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

var entityTypeToPB = map[db.EntityType]pb.EntityType{
	db.EntityTypeUnknown:   pb.EntityType_UNKNOWN,
	db.EntityTypePost:      pb.EntityType_POST,
	db.EntityTypeMap:       pb.EntityType_MAP,
	db.EntityTypeLocation:  pb.EntityType_LOCATION,
	db.EntityTypeCharacter: pb.EntityType_CHARACTER,
	db.EntityTypeImage:     pb.EntityType_IMAGE,
}

var entityTypeToDB = map[pb.EntityType]db.EntityType{
	pb.EntityType_UNKNOWN:   db.EntityTypeUnknown,
	pb.EntityType_POST:      db.EntityTypePost,
	pb.EntityType_MAP:       db.EntityTypeMap,
	pb.EntityType_LOCATION:  db.EntityTypeLocation,
	pb.EntityType_CHARACTER: db.EntityTypeCharacter,
	pb.EntityType_IMAGE:     db.EntityTypeImage,
}

func ConvertEntityTypeToPB(shape db.EntityType) pb.EntityType {
	if val, ok := entityTypeToPB[shape]; ok {
		return val
	}
	return pb.EntityType_UNKNOWN
}

func ConvertEntityTypeToDB(shape pb.EntityType) db.EntityType {
	if val, ok := entityTypeToDB[shape]; ok {
		return val
	}
	return db.EntityTypeUnknown
}
