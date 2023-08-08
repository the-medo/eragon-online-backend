package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertViewUser(user db.ViewUser) *pb.ViewUser {
	pbUser := &pb.ViewUser{
		Id:                user.ID,
		Username:          user.Username,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
		IsEmailVerified:   user.IsEmailVerified,
	}

	if user.ImgID.Valid == true {
		pbUser.ImgId = &user.ImgID.Int32
	}

	if user.AvatarImageUrl.Valid == true {
		pbUser.AvatarImageUrl = &user.AvatarImageUrl.String
	}

	if user.AvatarImageGuid.Valid == true {
		avatarImageGuid := user.AvatarImageGuid.UUID.String()
		pbUser.AvatarImageGuid = &avatarImageGuid
	}

	if user.IntroductionPostID.Valid == true {
		pbUser.IntroductionPostId = &user.IntroductionPostID.Int32
	}

	if user.IntroductionPostDeletedAt.Valid == true {
		pbUser.IntroductionPostDeletedAt = timestamppb.New(user.IntroductionPostDeletedAt.Time)
	}

	return pbUser
}
