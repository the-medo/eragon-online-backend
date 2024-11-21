package quests

import (
	"context"
	"github.com/the-medo/talebound-backend/api/servicecore"
	"github.com/the-medo/talebound-backend/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *ServiceQuests) CreateQuestCharacter(ctx context.Context, request *pb.CreateQuestCharacterRequest) (*pb.QuestCharacter, error) {
	violations := validateCreateQuestCharacter(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	//only character owner can add it to quest
	_, _, err := server.CheckModuleTypePermissions(ctx, db.ModuleTypeCharacter, request.GetCharacterId(), &servicecore.ModulePermission{
		NeedsSuperAdmin: true,
	})
	if err != nil {
		return nil, err
	}

	quest, err := server.Store.GetQuestByID(ctx, request.GetQuestId())
	if err != nil {
		return nil, err
	}
	character, err := server.Store.GetQuestByID(ctx, request.GetCharacterId())
	if err != nil {
		return nil, err
	}
	if quest.SystemID != character.SystemID {
		return nil, status.Errorf(codes.Internal, "failed to create quest character: quest and character systems do not match")
	}
	if quest.WorldID != character.WorldID {
		return nil, status.Errorf(codes.Internal, "failed to create quest character: quest and character worlds do not match")
	}

	argQuestCharacter := db.CreateQuestCharacterParams{
		QuestID:            request.GetQuestId(),
		CharacterID:        request.GetCharacterId(),
		MotivationalLetter: request.GetMotivationalLetter(),
	}

	newQuestCharacter, err := server.Store.CreateQuestCharacter(ctx, argQuestCharacter)
	if err != nil {
		return nil, err
	}

	rsp := converters.ConvertQuestCharacter(newQuestCharacter)

	return rsp, nil
}

func validateCreateQuestCharacter(req *pb.CreateQuestCharacterRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUniversalId(req.GetQuestId()); err != nil {
		violations = append(violations, e.FieldViolation("quest_id", err))
	}

	if err := validator.ValidateUniversalId(req.GetCharacterId()); err != nil {
		violations = append(violations, e.FieldViolation("character_id", err))
	}
	return violations
}
