package quests

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

func (server *ServiceQuests) GetQuestCharacters(ctx context.Context, req *pb.GetQuestCharactersRequest) (*pb.GetQuestCharactersResponse, error) {
	violations := validateGetQuestCharacters(req)

	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	questCharacters, err := server.Store.GetQuestCharactersByQuestID(ctx, req.GetQuestId())
	if err != nil {
		return nil, err
	}

	rsp := &pb.GetQuestCharactersResponse{
		QuestCharacters: make([]*pb.QuestCharacter, len(questCharacters)),
	}

	characterIds := make([]int32, len(questCharacters))

	for i, qc := range questCharacters {
		rsp.QuestCharacters[i] = converters.ConvertQuestCharacter(qc)
		characterIds[i] = qc.CharacterID
	}

	fetchInterface := &apihelpers.FetchInterface{
		CharacterIds: characterIds,
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

func validateGetQuestCharacters(req *pb.GetQuestCharactersRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateUniversalId(req.GetQuestId()); err != nil {
		violations = append(violations, e.FieldViolation("quest_id", err))
	}

	return violations
}
