package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertSystem(system db.System) *pb.System {
	pbSystem := &pb.System{
		Id:               system.ID,
		Name:             system.Name,
		Public:           system.Public,
		CreatedAt:        timestamppb.New(system.CreatedAt),
		ShortDescription: system.ShortDescription,
		BasedOn:          system.BasedOn,
	}

	return pbSystem
}
