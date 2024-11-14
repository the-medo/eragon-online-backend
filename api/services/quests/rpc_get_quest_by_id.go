package quests

import (
	"context"
	"github.com/the-medo/talebound-backend/converters"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceQuests) GetQuestById(ctx context.Context, req *pb.GetQuestByIdRequest) (*pb.Quest, error) {
	violations := validateGetQuestById(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	quest, err := server.Store.GetQuestByID(ctx, req.QuestId)
	if err != nil {
		return nil, err
	}

	return converters.ConvertQuest(quest), nil
}

func validateGetQuestById(req *pb.GetQuestByIdRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateModuleId(req.GetQuestId()); err != nil {
		violations = append(violations, e.FieldViolation("quest_id", err))
	}

	return violations
}
