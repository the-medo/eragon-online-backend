package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertModuleAdmin(dbModuleAdmin db.ModuleAdmin, dbUser db.User) *pb.ModuleAdmin {
	pbModuleAdmin := &pb.ModuleAdmin{
		ModuleId:           dbModuleAdmin.ModuleID,
		UserId:             dbModuleAdmin.UserID,
		User:               ConvertUser(dbUser),
		CreatedAt:          timestamppb.New(dbModuleAdmin.CreatedAt),
		SuperAdmin:         dbModuleAdmin.SuperAdmin,
		Approved:           dbModuleAdmin.Approved,
		MotivationalLetter: dbModuleAdmin.MotivationalLetter,
		AllowedEntityTypes: make([]pb.EntityType, len(dbModuleAdmin.AllowedEntityTypes)),
		AllowedMenu:        dbModuleAdmin.AllowedMenu,
	}

	for i, et := range dbModuleAdmin.AllowedEntityTypes {
		pbModuleAdmin.AllowedEntityTypes[i] = ConvertEntityTypeToPB(et)
	}

	return pbModuleAdmin
}
