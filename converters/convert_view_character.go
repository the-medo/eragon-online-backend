package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertViewCharacter(character db.ViewCharacter) *pb.ViewCharacter {
	pbViewCharacter := &pb.ViewCharacter{
		Id:               character.ID,
		ModuleId:         character.ModuleID,
		Name:             character.Name,
		Public:           character.Public,
		CreatedAt:        timestamppb.New(character.CreatedAt),
		ShortDescription: character.ShortDescription,
		Tags:             character.Tags,
		MenuId:           character.MenuID,
	}

	if character.WorldID.Valid {
		pbViewCharacter.WorldId = character.WorldID.Int32
	}

	if character.SystemID.Valid {
		pbViewCharacter.SystemId = character.SystemID.Int32
	}

	if character.HeaderImgID.Valid {
		pbViewCharacter.HeaderImgId = character.HeaderImgID.Int32
	}

	if character.ThumbnailImgID.Valid {
		pbViewCharacter.ThumbnailImgId = character.ThumbnailImgID.Int32
	}

	if character.AvatarImgID.Valid {
		pbViewCharacter.AvatarImgId = character.AvatarImgID.Int32
	}

	return pbViewCharacter
}
