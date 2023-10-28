package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

var moduleTypeToPB = map[db.ModuleType]pb.ModuleType{
	db.ModuleTypeUnknown:   pb.ModuleType_MODULE_TYPE_UNKNOWN,
	db.ModuleTypeWorld:     pb.ModuleType_MODULE_TYPE_WORLD,
	db.ModuleTypeQuest:     pb.ModuleType_MODULE_TYPE_QUEST,
	db.ModuleTypeSystem:    pb.ModuleType_MODULE_TYPE_SYSTEM,
	db.ModuleTypeCharacter: pb.ModuleType_MODULE_TYPE_CHARACTER,
}

var moduleTypeToDB = map[pb.ModuleType]db.ModuleType{
	pb.ModuleType_MODULE_TYPE_UNKNOWN:   db.ModuleTypeUnknown,
	pb.ModuleType_MODULE_TYPE_WORLD:     db.ModuleTypeWorld,
	pb.ModuleType_MODULE_TYPE_QUEST:     db.ModuleTypeQuest,
	pb.ModuleType_MODULE_TYPE_SYSTEM:    db.ModuleTypeSystem,
	pb.ModuleType_MODULE_TYPE_CHARACTER: db.ModuleTypeCharacter,
}

func ConvertModuleTypeToPB(shape db.ModuleType) pb.ModuleType {
	if val, ok := moduleTypeToPB[shape]; ok {
		return val
	}
	return pb.ModuleType_MODULE_TYPE_UNKNOWN
}

func ConvertModuleTypeToDB(shape pb.ModuleType) db.ModuleType {
	if val, ok := moduleTypeToDB[shape]; ok {
		return val
	}
	return db.ModuleTypeUnknown
}
