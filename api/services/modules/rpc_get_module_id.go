package modules

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/the-medo/talebound-backend/api/apihelpers"
	"github.com/the-medo/talebound-backend/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (server *ServiceModules) GetModuleId(ctx context.Context, req *pb.GetModuleIdRequest) (*pb.GetModuleIdResponse, error) {
	violations := validateGetModuleIdRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	module, err := server.Store.GetModule(ctx, db.GetModuleParams{
		WorldID: sql.NullInt32{
			Int32: req.GetWorldId(),
			Valid: req.WorldId != nil,
		},
		QuestID: sql.NullInt32{
			Int32: req.GetQuestId(),
			Valid: req.QuestId != nil,
		},
		CharacterID: sql.NullInt32{
			Int32: req.GetCharacterId(),
			Valid: req.CharacterId != nil,
		},
		SystemID: sql.NullInt32{
			Int32: req.GetSystemId(),
			Valid: req.SystemId != nil,
		},
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get module: %v", err)
	}

	rsp := &pb.GetModuleIdResponse{
		ModuleId:   module.ID,
		ModuleType: converters.ConvertModuleTypeToPB(module.ModuleType),
	}

	fetchInterface := &apihelpers.FetchInterface{
		ModuleIds: []int32{module.ID},
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

func validateGetModuleIdRequest(req *pb.GetModuleIdRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateModuleExtended(req.WorldId, req.QuestId, req.SystemId, req.CharacterId); err != nil {
		violations = append(violations, e.FieldViolation("modules", err))
	}
	return violations
}
