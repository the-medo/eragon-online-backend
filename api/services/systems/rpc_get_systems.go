package systems

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

func (server *ServiceSystems) GetSystems(ctx context.Context, req *pb.GetSystemsRequest) (*pb.GetSystemsResponse, error) {
	violations := validateGetSystems(req)

	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	limit, offset := apihelpers.GetDefaultQueryBoundaries(req.GetLimit(), req.GetOffset())

	arg := db.GetSystemsParams{
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

	systems, err := server.Store.GetSystems(ctx, arg)
	if err != nil {
		return nil, err
	}

	countArg := db.GetSystemsCountParams{
		IsPublic: arg.IsPublic,
		Tags:     req.GetTags(),
	}

	totalCount, err := server.Store.GetSystemsCount(ctx, countArg)
	if err != nil {
		return nil, err
	}

	moduleIds := make([]int32, len(systems))
	rsp := &pb.GetSystemsResponse{
		SystemIds:  make([]int32, len(systems)),
		TotalCount: int32(totalCount),
	}

	for i, system := range systems {
		rsp.SystemIds[i] = system.ID
		moduleIds[i] = system.ModuleID
	}

	fetchInterface := &apihelpers.FetchInterface{
		ModuleIds: moduleIds,
		SystemIds: rsp.SystemIds,
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

func validateGetSystems(req *pb.GetSystemsRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	fields := []string{"name", "created_at", "short_description", "activity_post_count", "activity_quest_count", "activity_resource_count"}

	if req.OrderBy != nil {
		if !validator.StringInSlice(req.GetOrderBy(), fields) {
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
