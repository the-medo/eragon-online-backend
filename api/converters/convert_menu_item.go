package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

func ConvertMenuItem(menuItem db.MenuItem) *pb.MenuItem {
	pbMenuItem := &pb.MenuItem{
		Id:       menuItem.ID,
		MenuId:   menuItem.MenuID,
		Code:     menuItem.MenuItemCode,
		Name:     menuItem.Name,
		Position: menuItem.Position,
	}

	if menuItem.ParentItemID.Valid == true {
		pbMenuItem.ParentItemId = &menuItem.ParentItemID.Int32
	}

	if menuItem.DescriptionPostID.Valid == true {
		pbMenuItem.DescriptionPostId = &menuItem.DescriptionPostID.Int32
	}

	return pbMenuItem
}
