package quests

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
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

func (server *ServiceQuests) GetQuests(ctx context.Context, req *pb.GetQuestsRequest) (*pb.GetQuestsResponse, error) {
	violations := validateGetQuests(req)

	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	limit, offset := apihelpers.GetDefaultQueryBoundaries(req.GetLimit(), req.GetOffset())

	arg := db.GetQuestsParams{
		PageLimit:  sql.NullInt32{Int32: limit, Valid: true},
		PageOffset: sql.NullInt32{Int32: offset, Valid: true},
		Tags:       req.GetTags(),
		OrderBy:    sql.NullString{String: "created_at", Valid: true},
		WorldID: sql.NullInt32{
			Int32: req.GetWorldId(),
			Valid: req.WorldId != nil,
		},
		SystemID: sql.NullInt32{
			Int32: req.GetSystemId(),
			Valid: req.SystemId != nil,
		},
		CanJoin: sql.NullBool{
			Bool:  req.GetCanJoin(),
			Valid: req.CanJoin != nil,
		},
		Status: db.NullQuestStatus{
			QuestStatus: converters.ConvertQuestStatusToDB(req.GetStatus()),
			Valid:       req.Status != nil,
		},
	}

	countArg := db.GetQuestsCountParams{
		IsPublic: sql.NullBool{},
		Tags:     req.GetTags(),
		WorldID:  arg.WorldID,
		SystemID: arg.SystemID,
		CanJoin:  arg.CanJoin,
		Status:   arg.Status,
	}

	if req.Public != nil {
		arg.IsPublic = sql.NullBool{Bool: req.GetPublic(), Valid: true}
		countArg.IsPublic = sql.NullBool{
			Bool:  req.GetPublic(),
			Valid: true,
		}
	}

	if req.OrderBy != nil {
		arg.OrderBy.String = req.GetOrderBy()
	}

	quests, err := server.Store.GetQuests(ctx, arg)
	if err != nil {
		return nil, err
	}

	totalCount, err := server.Store.GetQuestsCount(ctx, countArg)
	if err != nil {
		return nil, err
	}

	moduleIds := make([]int32, len(quests))
	rsp := &pb.GetQuestsResponse{
		QuestIds:   make([]int32, len(quests)),
		TotalCount: int32(totalCount),
	}

	for i, quest := range quests {
		rsp.QuestIds[i] = quest.ID
		moduleIds[i] = quest.ModuleID
	}

	fetchInterface := &apihelpers.FetchInterface{
		ModuleIds: moduleIds,
		QuestIds:  rsp.QuestIds,
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

func validateGetQuests(req *pb.GetQuestsRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	//util.CleanEmptyStrings(&req.Tags)

	fields := []string{"name", "created_at", "short_description"}

	if req.OrderBy != nil {
		if validator.StringInSlice(req.GetOrderBy(), fields) == false {
			violations = append(violations, e.FieldViolation("order_by", fmt.Errorf("invalid field to order by %s", req.GetOrderBy())))
		}
	}

	if req.Limit != nil {
		if err := validator.ValidateLimit(req.GetLimit()); err != nil {
			violations = append(violations, e.FieldViolation("limit", err))
		}
	}

	if req.Offset != nil {
		if err := validator.ValidateOffset(req.GetOffset()); err != nil {
			violations = append(violations, e.FieldViolation("offset", err))
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
