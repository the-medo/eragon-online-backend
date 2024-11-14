package characters

import (
	"context"
	"github.com/the-medo/talebound-backend/converters"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceCharacters) GetCharacterById(ctx context.Context, req *pb.GetCharacterByIdRequest) (*pb.Character, error) {
	violations := validateGetCharacterById(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	character, err := server.Store.GetCharacterByID(ctx, req.CharacterId)
	if err != nil {
		return nil, err
	}

	return converters.ConvertCharacter(character), nil
}

func validateGetCharacterById(req *pb.GetCharacterByIdRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateModuleId(req.GetCharacterId()); err != nil {
		violations = append(violations, e.FieldViolation("character_id", err))
	}

	return violations
}
