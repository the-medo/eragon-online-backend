package quests

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/api/servicecore"
	"github.com/the-medo/talebound-backend/constants"
	"github.com/the-medo/talebound-backend/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *ServiceQuests) UpdateQuest(ctx context.Context, req *pb.UpdateQuestRequest) (*pb.Quest, error) {
	violations := validateUpdateQuestRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	var needsEntityPermission []db.EntityType

	_, _, err := server.CheckModuleTypePermissions(ctx, db.ModuleTypeQuest, req.GetQuestId(), &servicecore.ModulePermission{
		NeedsSuperAdmin:       true,
		NeedsEntityPermission: &needsEntityPermission,
	})

	if err != nil {
		return nil, err
	}

	// WorldID and SystemID are updateable only when they are "Universal world" or "Universal system"
	// otherwise, another related entities can be out of sync with what is possible
	if req.WorldId != nil || req.SystemId != nil {
		questCharacters, err := server.Store.GetQuestCharactersByQuestID(ctx, req.GetQuestId())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update quest: %s", err)
		}

		if len(questCharacters) > 0 {
			return nil, status.Errorf(codes.Internal, "failed to update quest: can not change world or system when characters are in a quest")
		}

		if req.WorldId != nil {
			world, err := server.Store.GetWorldByID(ctx, req.GetWorldId())
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to update quest: %s", err)
			}
			if world.Name != constants.UniversalWorldName {
				return nil, status.Errorf(codes.Internal, "failed to update quest: world can be changed only from world named %s", constants.UniversalWorldName)
			}
		}

		if req.SystemId != nil {
			system, err := server.Store.GetSystemByID(ctx, req.GetSystemId())
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to update quest: %s", err)
			}
			if system.Name != constants.UniversalSystemName {
				return nil, status.Errorf(codes.Internal, "failed to update quest: system can be changed only from system named %s", constants.UniversalSystemName)
			}
		}
	}

	arg := db.UpdateQuestParams{
		QuestID:          req.GetQuestId(),
		Name:             sql.NullString{String: req.GetName(), Valid: req.Name != nil},
		ShortDescription: sql.NullString{String: req.GetShortDescription(), Valid: req.ShortDescription != nil},
		WorldID:          sql.NullInt32{Int32: req.GetWorldId(), Valid: req.WorldId != nil},
		SystemID:         sql.NullInt32{Int32: req.GetSystemId(), Valid: req.SystemId != nil},
		CanJoin:          sql.NullBool{Bool: req.GetCanJoin(), Valid: req.CanJoin != nil},
		Status:           db.NullQuestStatus{QuestStatus: converters.ConvertQuestStatusToDB(req.GetStatus()), Valid: req.Status != nil},
	}

	quest, err := server.Store.UpdateQuest(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update quest: %s", err)
	}

	return converters.ConvertQuest(quest), nil
}

func validateUpdateQuestRequest(req *pb.UpdateQuestRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateModuleId(req.GetQuestId()); err != nil {
		violations = append(violations, e.FieldViolation("quest_id", err))
	}

	if req.Name != nil {
		if err := validator.ValidateModuleName(req.GetName()); err != nil {
			violations = append(violations, e.FieldViolation("name", err))
		}
	}

	if req.ShortDescription != nil {
		if err := validator.ValidateModuleShortDescription(req.GetShortDescription()); err != nil {
			violations = append(violations, e.FieldViolation("short_description", err))
		}
	}

	if req.WorldId != nil {
		if err := validator.ValidateUniversalId(req.GetWorldId()); err != nil {
			violations = append(violations, e.FieldViolation("world_id", err))
		}
	}

	if req.SystemId != nil {
		if err := validator.ValidateUniversalId(req.GetSystemId()); err != nil {
			violations = append(violations, e.FieldViolation("system_id", err))
		}
	}

	return violations
}
