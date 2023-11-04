package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertGetUserRow(user db.GetUsersRow) *pb.User {
	pbUser := &pb.User{
		Id:                user.ID,
		Username:          user.Username,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
		IsEmailVerified:   user.IsEmailVerified,
	}

	if user.ImgID.Valid == true {
		pbUser.ImgId = &user.ImgID.Int32
		pbUser.Img = &pb.Image{
			Id:          user.ImgID.Int32,
			ImgGuid:     user.ImgGuid.UUID.String(),
			ImageTypeId: user.ImageTypeID.Int32,
			Name:        user.Name.String,
			Url:         user.Url.String,
			BaseUrl:     user.BaseUrl.String,
			CreatedAt:   timestamppb.New(user.CreatedAt_2.Time),
		}
	}

	return pbUser
}
