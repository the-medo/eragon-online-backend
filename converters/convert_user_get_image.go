package converters

import (
	"context"
	"github.com/the-medo/talebound-backend/api/servicecore"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

func ConvertUserGetImage(server *servicecore.ServiceCore, ctx context.Context, user db.User) *pb.User {
	pbUser := ConvertUser(user, nil)

	if user.ImgID.Valid == true {
		img, err := server.Store.GetImageById(ctx, *pbUser.ImgId)
		if err != nil {
			return nil
		}
		pbUser.Img = ConvertImage(&img)
	}

	return pbUser
}
