package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertViewWorld(world db.ViewWorld) *pb.ViewWorld {
	pbViewWorld := &pb.ViewWorld{
		Id:               world.ID,
		ModuleId:         world.ModuleID,
		Name:             world.Name,
		Public:           world.Public,
		CreatedAt:        timestamppb.New(world.CreatedAt),
		ShortDescription: world.ShortDescription,
		BasedOn:          world.BasedOn,
		Tags:             world.Tags,
		MenuId:           world.MenuID,
	}

	if world.HeaderImgID.Valid {
		pbViewWorld.HeaderImgId = world.HeaderImgID.Int32
	}

	if world.ThumbnailImgID.Valid {
		pbViewWorld.ThumbnailImgId = world.ThumbnailImgID.Int32
	}

	if world.AvatarImgID.Valid {
		pbViewWorld.AvatarImgId = world.AvatarImgID.Int32
	}

	return pbViewWorld
}
