package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertWorld(world db.World) *pb.World {
	pbWorld := &pb.World{
		Id:               world.ID,
		Name:             world.Name,
		Public:           world.Public,
		CreatedAt:        timestamppb.New(world.CreatedAt),
		ShortDescription: world.ShortDescription,
		BasedOn:          world.BasedOn,
	}

	return pbWorld
}
