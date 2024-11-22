package quests

import (
	"context"
	"github.com/the-medo/talebound-backend/api/servicecore"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *ServiceQuests) DeleteQuestCharacter(ctx context.Context, req *pb.DeleteQuestCharacterRequest) (*emptypb.Empty, error) {
	violations := validateDeleteQuestCharacterRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	//Only quest or character owner can delete it
	_, _, err := server.CheckModuleTypePermissions(ctx, db.ModuleTypeQuest, req.GetQuestId(), &servicecore.ModulePermission{
		NeedsSuperAdmin: true,
	})
	if err != nil {
		_, _, err := server.CheckModuleTypePermissions(ctx, db.ModuleTypeCharacter, req.GetCharacterId(), &servicecore.ModulePermission{
			NeedsSuperAdmin: true,
		})
		if err != nil {
			return nil, err
		}
	}

	arg := db.DeleteQuestCharacterParams{
		QuestID:     req.GetQuestId(),
		CharacterID: req.GetCharacterId(),
	}

	err = server.Store.DeleteQuestCharacter(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete quest character: %s", err)
	}

	return &emptypb.Empty{}, nil
}

func validateDeleteQuestCharacterRequest(req *pb.DeleteQuestCharacterRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateModuleId(req.GetQuestId()); err != nil {
		violations = append(violations, e.FieldViolation("quest_id", err))
	}

	if err := validator.ValidateModuleId(req.GetCharacterId()); err != nil {
		violations = append(violations, e.FieldViolation("character_id", err))
	}

	return violations
}
