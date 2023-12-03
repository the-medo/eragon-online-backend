package users

import (
	"context"
	"github.com/goccy/go-json"
	"github.com/the-medo/talebound-backend/api/apihelpers"
	"github.com/the-medo/talebound-backend/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

	moduleIDs := make([]int32, len(userModules))
	worldIDs := make([]int32, 0)
	systemIDs := make([]int32, 0)
	characterIDs := make([]int32, 0)
	questIDs := make([]int32, 0)

	rsp := &pb.GetUserModulesResponse{
		UserModules: make([]*pb.UserModule, len(userModules)),
	}

	for i, userModule := range userModules {

		if userModule.ModuleType == db.ModuleTypeWorld {
			worldIDs = append(worldIDs, userModule.WorldID.Int32)
		} else if userModule.ModuleType == db.ModuleTypeSystem {
			systemIDs = append(systemIDs, userModule.SystemID.Int32)
		} else if userModule.ModuleType == db.ModuleTypeCharacter {
			characterIDs = append(characterIDs, userModule.CharacterID.Int32)
		} else if userModule.ModuleType == db.ModuleTypeQuest {
			questIDs = append(questIDs, userModule.QuestID.Int32)
		}

		moduleIDs[i] = userModule.ID
		rsp.UserModules[i] = converters.ConvertGetUserModulesRow(userModule)
	}

	fetchInterface := &apihelpers.FetchInterface{
		ModuleIds:    moduleIDs,
		WorldIds:     worldIDs,
		SystemIds:    systemIDs,
		CharacterIds: characterIDs,
		QuestIds:     questIDs,
	}

	fetchIdsHeader, err := json.Marshal(fetchInterface)

	md := metadata.Pairs(
		"X-Fetch-Ids", string(fetchIdsHeader),
	)
	err = grpc.SendHeader(ctx, md)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func validateGetUserModules(req *pb.GetUserModulesRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUserId(req.GetUserId()); err != nil {
		violations = append(violations, e.FieldViolation("user_id", err))
	}

	return violations
}
