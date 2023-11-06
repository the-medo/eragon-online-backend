package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

func ConvertMenuItemPost(menuItemPost db.MenuItemPost, post db.ViewPost) *pb.MenuItemPost {
	pbMenuItemPost := &pb.MenuItemPost{
		MenuItemId: menuItemPost.MenuItemID.Int32,
		PostId:     menuItemPost.PostID,
		Position:   menuItemPost.Position,
		Post:       ConvertViewPost(post),
	}

	return pbMenuItemPost
}
