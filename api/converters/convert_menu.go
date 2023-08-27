package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

func ConvertViewMenu(menu db.ViewMenu) *pb.ViewMenu {
	pbMenu := &pb.ViewMenu{
		Id:   menu.ID,
		Code: menu.MenuCode,
	}

	if menu.MenuHeaderImgID.Valid == true {
		pbMenu.HeaderImageId = &menu.MenuHeaderImgID.Int32
	}

	if menu.HeaderImageUrl.Valid == true {
		pbMenu.HeaderImageUrl = &menu.HeaderImageUrl.String
	}

	return pbMenu
}
