package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertWorldAdminRow(dbWorldAdmin db.GetWorldAdminsRow) *pb.WorldAdmin {
	pbViewUser := &pb.ViewUser{
		Id:                dbWorldAdmin.ID,
		Username:          dbWorldAdmin.Username,
		Email:             dbWorldAdmin.Email,
		PasswordChangedAt: timestamppb.New(dbWorldAdmin.PasswordChangedAt),
		CreatedAt:         timestamppb.New(dbWorldAdmin.CreatedAt),
		IsEmailVerified:   dbWorldAdmin.IsEmailVerified,
	}

	if dbWorldAdmin.ImgID.Valid == true {
		pbViewUser.ImgId = &dbWorldAdmin.ImgID.Int32
	}

	if dbWorldAdmin.AvatarImageUrl.Valid == true {
		pbViewUser.AvatarImageUrl = &dbWorldAdmin.AvatarImageUrl.String
	}

	if dbWorldAdmin.AvatarImageGuid.Valid == true {
		avatarImageGuid := dbWorldAdmin.AvatarImageGuid.UUID.String()
		pbViewUser.AvatarImageGuid = &avatarImageGuid
	}

	if dbWorldAdmin.IntroductionPostID.Valid == true {
		pbViewUser.IntroductionPostId = &dbWorldAdmin.IntroductionPostID.Int32
	}

	if dbWorldAdmin.IntroductionPostDeletedAt.Valid == true {
		pbViewUser.IntroductionPostDeletedAt = timestamppb.New(dbWorldAdmin.IntroductionPostDeletedAt.Time)
	}

	pbWorldAdmin := &pb.WorldAdmin{
		WorldId:            dbWorldAdmin.WorldID,
		UserId:             dbWorldAdmin.ID,
		User:               pbViewUser,
		CreatedAt:          timestamppb.New(dbWorldAdmin.WorldAdminCreatedAt),
		SuperAdmin:         dbWorldAdmin.WorldAdminSuperAdmin,
		Approved:           dbWorldAdmin.WorldAdminApproved,
		MotivationalLetter: dbWorldAdmin.WorldAdminMotivationalLetter,
	}

	return pbWorldAdmin
}
