package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertModuleAdminRow(dbModuleAdmin db.GetModuleAdminsRow) *pb.ModuleAdmin {
	pbUser := &pb.User{
		Id:                dbModuleAdmin.ID,
		Username:          dbModuleAdmin.Username,
		Email:             dbModuleAdmin.Email,
		PasswordChangedAt: timestamppb.New(dbModuleAdmin.PasswordChangedAt),
		CreatedAt:         timestamppb.New(dbModuleAdmin.CreatedAt),
		IsEmailVerified:   dbModuleAdmin.IsEmailVerified,
	}

	if dbModuleAdmin.ImgID.Valid == true {
		pbUser.ImgId = &dbModuleAdmin.ImgID.Int32
	}

	if dbModuleAdmin.IntroductionPostID.Valid == true {
		pbUser.IntroductionPostId = &dbModuleAdmin.IntroductionPostID.Int32
	}

	pbModuleAdmin := &pb.ModuleAdmin{
		ModuleId:           dbModuleAdmin.ModuleID,
		UserId:             dbModuleAdmin.ID,
		User:               pbUser,
		CreatedAt:          timestamppb.New(dbModuleAdmin.ModuleAdminCreatedAt),
		SuperAdmin:         dbModuleAdmin.ModuleAdminSuperAdmin,
		Approved:           dbModuleAdmin.ModuleAdminApproved,
		MotivationalLetter: dbModuleAdmin.ModuleAdminMotivationalLetter,
		AllowedEntityTypes: make([]pb.EntityType, len(dbModuleAdmin.ModuleAdminAllowedEntityTypes)),
		AllowedMenu:        dbModuleAdmin.ModuleAdminAllowedMenu,
	}

	for i, entityType := range dbModuleAdmin.ModuleAdminAllowedEntityTypes {
		pbModuleAdmin.AllowedEntityTypes[i] = ConvertEntityTypeToPB(entityType)
	}

	return pbModuleAdmin
}
