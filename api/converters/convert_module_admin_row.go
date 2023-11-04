package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertModuleAdminRow(dbModuleAdmin db.GetModuleAdminsRow) *pb.ModuleAdmin {
	pbViewUser := &pb.ViewUser{
		Id:                dbModuleAdmin.ID,
		Username:          dbModuleAdmin.Username,
		Email:             dbModuleAdmin.Email,
		PasswordChangedAt: timestamppb.New(dbModuleAdmin.PasswordChangedAt),
		CreatedAt:         timestamppb.New(dbModuleAdmin.CreatedAt),
		IsEmailVerified:   dbModuleAdmin.IsEmailVerified,
	}

	if dbModuleAdmin.ImgID.Valid == true {
		pbViewUser.ImgId = &dbModuleAdmin.ImgID.Int32
	}

	if dbModuleAdmin.AvatarImageUrl.Valid == true {
		pbViewUser.AvatarImageUrl = &dbModuleAdmin.AvatarImageUrl.String
	}

	if dbModuleAdmin.AvatarImageGuid.Valid == true {
		avatarImageGuid := dbModuleAdmin.AvatarImageGuid.UUID.String()
		pbViewUser.AvatarImageGuid = &avatarImageGuid
	}

	if dbModuleAdmin.IntroductionPostID.Valid == true {
		pbViewUser.IntroductionPostId = &dbModuleAdmin.IntroductionPostID.Int32
	}

	if dbModuleAdmin.IntroductionPostDeletedAt.Valid == true {
		pbViewUser.IntroductionPostDeletedAt = timestamppb.New(dbModuleAdmin.IntroductionPostDeletedAt.Time)
	}

	pbModuleAdmin := &pb.ModuleAdmin{
		WorldId:            dbModuleAdmin.ModuleID,
		UserId:             dbModuleAdmin.ID,
		User:               pbViewUser,
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
