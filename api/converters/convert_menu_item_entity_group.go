package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

func ConvertMenuItemEntityGroup(item db.MenuItemEntityGroup) *pb.MenuItemEntityGroup {
	pbItem := &pb.MenuItemEntityGroup{
		MenuId:        item.MenuID,
		EntityGroupId: item.EntityGroupID,
		Position:      item.Position,
	}

	if item.MenuItemID.Valid {
		pbItem.MenuItemId = item.MenuItemID.Int32
	}

	return pbItem
}
