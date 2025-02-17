package worlds

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

func (server *ServiceWorlds) GetWorlds(ctx context.Context, req *pb.GetWorldsRequest) (*pb.GetWorldsResponse, error) {
	violations := validateGetWorlds(req)

	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	limit, offset := apihelpers.GetDefaultQueryBoundaries(req.GetLimit(), req.GetOffset())

	arg := db.GetWorldsParams{
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

	worlds, err := server.Store.GetWorlds(ctx, arg)
	if err != nil {
		return nil, err
	}

	countArg := db.GetWorldsCountParams{
		IsPublic: arg.IsPublic,
		Tags:     req.GetTags(),
	}

	totalCount, err := server.Store.GetWorldsCount(ctx, countArg)
	if err != nil {
		return nil, err
	}

	moduleIds := make([]int32, len(worlds))
	rsp := &pb.GetWorldsResponse{
		WorldIds:   make([]int32, len(worlds)),
		TotalCount: int32(totalCount),
	}

	for i, world := range worlds {
		rsp.WorldIds[i] = world.ID
		moduleIds[i] = world.ModuleID
	}

	fetchInterface := &apihelpers.FetchInterface{
		ModuleIds: moduleIds,
		WorldIds:  rsp.WorldIds,
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

func validateGetWorlds(req *pb.GetWorldsRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	//util.CleanEmptyStrings(&req.Tags)

	fields := []string{"name", "created_at", "short_description", "activity_post_count", "activity_quest_count", "activity_resource_count"}

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

	return violations
}
