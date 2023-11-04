package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertViewWorld(world db.ViewWorld) *pb.World {
	pbWorld := &pb.World{
		Id:               world.ID,
		Name:             world.Name,
		Public:           world.Public,
		CreatedAt:        timestamppb.New(world.CreatedAt),
		ShortDescription: world.ShortDescription,
		BasedOn:          world.BasedOn,
		ImageAvatar:      world.ImageAvatar.String,
		ImageThumbnail:   world.ImageThumbnail.String,
		ImageHeader:      world.ImageHeader.String,
		Tags:             world.Tags,
	}

	if world.DescriptionPostID.Valid == true {
		pbWorld.DescriptionPostId = &world.DescriptionPostID.Int32
	}

	if world.MenuID.Valid == true {
		pbWorld.WorldMenuId = world.MenuID.Int32
	}

	return pbWorld
}
