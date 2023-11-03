package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertModuleAdmin(dbModuleAdmin db.ModuleAdmin, dbViewUser db.ViewUser) *pb.ModuleAdmin {
	pbModuleAdmin := &pb.ModuleAdmin{
		WorldId:            dbModuleAdmin.ModuleID,
		UserId:             dbModuleAdmin.UserID,
		User:               ConvertViewUser(dbViewUser),
		CreatedAt:          timestamppb.New(dbModuleAdmin.CreatedAt),
		SuperAdmin:         dbModuleAdmin.SuperAdmin,
		Approved:           dbModuleAdmin.Approved,
		MotivationalLetter: dbModuleAdmin.MotivationalLetter,
	}

	return pbModuleAdmin
}
