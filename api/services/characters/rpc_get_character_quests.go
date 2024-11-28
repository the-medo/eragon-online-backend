package characters

import (
	"context"
	"encoding/json"
	"github.com/the-medo/talebound-backend/api/apihelpers"
	"github.com/the-medo/talebound-backend/converters"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (server *ServiceCharacters) GetCharacterQuests(ctx context.Context, req *pb.GetCharacterQuestsRequest) (*pb.GetQuestCharactersResponse, error) {
	violations := validateGetCharacterQuests(req)

	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	questCharacters, err := server.Store.GetQuestCharactersByCharacterID(ctx, req.GetCharacterId())
	if err != nil {
		return nil, err
	}

	rsp := &pb.GetQuestCharactersResponse{
		QuestCharacters: make([]*pb.QuestCharacter, len(questCharacters)),
	}

	questIds := make([]int32, len(questCharacters))

	for i, qc := range questCharacters {
		rsp.QuestCharacters[i] = converters.ConvertQuestCharacter(qc)
		questIds[i] = qc.QuestID
	}

	fetchInterface := &apihelpers.FetchInterface{
		QuestIds: questIds,
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

func validateGetCharacterQuests(req *pb.GetCharacterQuestsRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateUniversalId(req.GetCharacterId()); err != nil {
		violations = append(violations, e.FieldViolation("character_id", err))
	}

	return violations
}
