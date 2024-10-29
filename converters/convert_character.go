package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertCharacter(character db.Character) *pb.Character {
	pbCharacter := &pb.Character{
		Id:               character.ID,
		Name:             character.Name,
		Public:           character.Public,
		CreatedAt:        timestamppb.New(character.CreatedAt),
		ShortDescription: character.ShortDescription,
	}

	if character.WorldID.Valid {
		pbCharacter.WorldId = character.WorldID.Int32
	}

	if character.SystemID.Valid {
		pbCharacter.SystemId = character.SystemID.Int32
	}

	return pbCharacter
}
