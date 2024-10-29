package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertViewSystem(system db.ViewSystem) *pb.ViewSystem {
	pbViewSystem := &pb.ViewSystem{
		Id:               system.ID,
		ModuleId:         system.ModuleID,
		Name:             system.Name,
		Public:           system.Public,
		CreatedAt:        timestamppb.New(system.CreatedAt),
		ShortDescription: system.ShortDescription,
		BasedOn:          system.BasedOn,
		Tags:             system.Tags,
		MenuId:           system.MenuID,
	}

	if system.HeaderImgID.Valid {
		pbViewSystem.HeaderImgId = system.HeaderImgID.Int32
	}

	if system.ThumbnailImgID.Valid {
		pbViewSystem.ThumbnailImgId = system.ThumbnailImgID.Int32
	}

	if system.AvatarImgID.Valid {
		pbViewSystem.AvatarImgId = system.AvatarImgID.Int32
	}

	return pbViewSystem
}
