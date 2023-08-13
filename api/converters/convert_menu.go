package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

func ConvertMenu(menu db.Menu) *pb.Menu {
	pbMenu := &pb.Menu{
		Id:   menu.ID,
		Code: menu.MenuCode,
	}

	if menu.MenuHeaderImgID.Valid == true {
		pbMenu.HeaderImageId = &menu.MenuHeaderImgID.Int32
	}

	return pbMenu
}
