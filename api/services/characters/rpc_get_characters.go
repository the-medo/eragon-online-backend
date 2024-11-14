package characters

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/the-medo/talebound-backend/api/apihelpers"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (server *ServiceCharacters) GetCharacters(ctx context.Context, req *pb.GetCharactersRequest) (*pb.GetCharactersResponse, error) {
	violations := validateGetCharacters(req)

	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	limit, offset := apihelpers.GetDefaultQueryBoundaries(req.GetLimit(), req.GetOffset())

	arg := db.GetCharactersParams{
		PageLimit:  limit,
		PageOffset: offset,
		Tags:       req.GetTags(),
		OrderBy:    "created_at",
	}

	if req.Public != nil {
		arg.IsPublic = req.GetPublic()
	} else {
		arg.IsPublic = true
	}

	if req.OrderBy != nil {
		arg.OrderBy = req.GetOrderBy()
	}

	characters, err := server.Store.GetCharacters(ctx, arg)
	if err != nil {
		return nil, err
	}

	countArg := db.GetCharactersCountParams{
		IsPublic: arg.IsPublic,
		Tags:     req.GetTags(),
	}

	totalCount, err := server.Store.GetCharactersCount(ctx, countArg)
	if err != nil {
		return nil, err
	}

	moduleIds := make([]int32, len(characters))
	rsp := &pb.GetCharactersResponse{
		CharacterIds: make([]int32, len(characters)),
		TotalCount:   int32(totalCount),
	}

	for i, character := range characters {
		rsp.CharacterIds[i] = character.ID
		moduleIds[i] = character.ModuleID
	}

	fetchInterface := &apihelpers.FetchInterface{
		ModuleIds:    moduleIds,
		CharacterIds: rsp.CharacterIds,
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

func validateGetCharacters(req *pb.GetCharactersRequest) (violations []*errdetails.BadRequest_FieldViolation) {
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
