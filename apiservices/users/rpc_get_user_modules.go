package users

import (
	"context"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/api/e"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceUsers) GetUserModules(ctx context.Context, req *pb.GetUserModulesRequest) (*pb.GetUserModulesResponse, error) {
	violations := validateGetUserModules(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	userModules, err := server.Store.GetUserModules(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	worldIDs := make([]int32, 0)
	systemIDs := make([]int32, 0)
	characterIDs := make([]int32, 0)
	questIDs := make([]int32, 0)

	for _, module := range userModules {
		if module.ModuleType == db.ModuleTypeWorld {
			worldIDs = append(worldIDs, module.WorldID.Int32)
		} else if module.ModuleType == db.ModuleTypeSystem {
			systemIDs = append(systemIDs, module.SystemID.Int32)
		} else if module.ModuleType == db.ModuleTypeCharacter {
			characterIDs = append(characterIDs, module.CharacterID.Int32)
		} else if module.ModuleType == db.ModuleTypeQuest {
			questIDs = append(questIDs, module.QuestID.Int32)
		}
	}

	rsp := &pb.GetUserModulesResponse{
		Modules: make([]*pb.ViewModule, len(userModules)),
		Worlds:  make([]*pb.World, len(worldIDs)),
	}

	if len(worldIDs) > 0 {
		worlds, err := server.Store.GetWorldsByIDs(ctx, worldIDs)
		if err != nil {
			return nil, err
		}
		for i, world := range worlds {
			rsp.Worlds[i] = converters.ConvertViewWorld(world)
		}
	}

	//TODO: Implement the rest of these

	/*
		if len(systemIDs) > 0 {
			systems, err := server.Store.GetSystemsByIDs(ctx, systemIDs)
			if err != nil {
				return nil, err
			}
			for i, system := range systems {
				rsp.Worlds[i] = converters.ConvertViewSystem(system)
			}
		}*/

	/*
		if len(characterIDs) > 0 {
			characters, err := server.Store.GetCharactersByIDs(ctx, characterIDs)
			if err != nil {
				return nil, err
			}
			for i, character := range characters {
				rsp.Worlds[i] = converters.ConvertViewCharacter(character)
			}
		}*/

	/*
		if len(questIDs) > 0 {
			quests, err := server.Store.GetQuestsByIDs(ctx, questIDs)
			if err != nil {
				return nil, err
			}
			for i, quest := range quests {
				rsp.Worlds[i] = converters.ConvertViewQuest(quest)
			}
		}*/

	return rsp, nil
}

func validateGetUserModules(req *pb.GetUserModulesRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUserId(req.GetUserId()); err != nil {
		violations = append(violations, e.FieldViolation("user_id", err))
	}

	return violations
}
