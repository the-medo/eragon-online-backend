package quests

import (
	"context"
	"database/sql"
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

func (server *ServiceQuests) UpdateQuestCharacter(ctx context.Context, req *pb.UpdateQuestCharacterRequest) (*pb.QuestCharacter, error) {
	violations := validateUpdateQuestCharacterRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	//Only quest owner can change Approved
	if req.Approved != nil {
		_, _, err := server.CheckModuleTypePermissions(ctx, db.ModuleTypeQuest, req.GetQuestId(), &servicecore.ModulePermission{
			NeedsSuperAdmin: true,
		})
		if err != nil {
			return nil, err
		}
	}

	// Only character owner can change MotivationalLetter
	if req.MotivationalLetter != nil {
		_, _, err := server.CheckModuleTypePermissions(ctx, db.ModuleTypeCharacter, req.GetCharacterId(), &servicecore.ModulePermission{
			NeedsSuperAdmin: true,
		})
		if err != nil {
			return nil, err
		}
	}

	arg := db.UpdateQuestCharacterParams{
		QuestID:     req.GetQuestId(),
		CharacterID: req.GetCharacterId(),
		Approved: sql.NullInt32{
			Int32: req.GetApproved(),
			Valid: req.Approved != nil,
		},
		MotivationalLetter: sql.NullString{
			String: req.GetMotivationalLetter(),
			Valid:  req.MotivationalLetter != nil,
		},
	}

	questCharacter, err := server.Store.UpdateQuestCharacter(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update quest character: %s", err)
	}

	return converters.ConvertQuestCharacter(questCharacter), nil
}

func validateUpdateQuestCharacterRequest(req *pb.UpdateQuestCharacterRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if req.MotivationalLetter == nil && req.Approved == nil {
		violations = append(violations, e.FieldViolation("motivational_letter", status.Errorf(codes.Internal, "failed to update quest character: no update field is present")))
	}

	if err := validator.ValidateModuleId(req.GetQuestId()); err != nil {
		violations = append(violations, e.FieldViolation("quest_id", err))
	}

	if err := validator.ValidateModuleId(req.GetCharacterId()); err != nil {
		violations = append(violations, e.FieldViolation("character_id", err))
	}

	if err := validator.ValidateInt(req.GetApproved(), 0, 2); err != nil {
		violations = append(violations, e.FieldViolation("approved", err))
	}

	return violations
}
