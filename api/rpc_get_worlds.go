package api

import (
	"context"
	"fmt"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *Server) GetWorlds(ctx context.Context, req *pb.GetWorldsRequest) (*pb.GetWorldsResponse, error) {
	violations := validateGetWorlds(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	limit, offset := GetDefaultQueryBoundaries(req.GetLimit(), req.GetOffset())

	arg := db.GetWorldsParams{
		PageLimit:  limit,
		PageOffset: offset,
	}

	if req.Public != nil {
		arg.IsPublic = req.GetPublic()
	} else {
		arg.IsPublic = true
	}

	if req.OrderBy != nil {
		arg.OrderResult = true
		arg.OrderBy = req.GetOrderBy()
	}

	worlds, err := server.store.GetWorlds(ctx, arg)
	if err != nil {
		return nil, err
	}

	totalCount, err := server.store.GetWorldsCount(ctx, arg.IsPublic)
	if err != nil {
		return nil, err
	}

	rsp := &pb.GetWorldsResponse{
		Worlds:     make([]*pb.World, len(worlds)),
		TotalCount: int32(totalCount),
	}

	for i, world := range worlds {
		rsp.Worlds[i] = convertWorld(world)
	}

	return rsp, nil
}

func validateGetWorlds(req *pb.GetWorldsRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	fields := []string{"name", "created_at", "short_description", "activity_post_count", "activity_quest_count", "activity_resource_count"}

	if req.OrderBy != nil {
		if validator.StringInSlice(req.GetOrderBy(), fields) == false {
			violations = append(violations, FieldViolation("order_by", fmt.Errorf("invalid field to order by %s", req.GetOrderBy())))
		}
	}

	if req.Limit != nil {
		if err := validator.ValidateLimit(req.GetLimit()); err != nil {
			violations = append(violations, FieldViolation("limit", err))
		}
	}

	if req.Offset != nil {
		if err := validator.ValidateOffset(req.GetOffset()); err != nil {
			violations = append(violations, FieldViolation("offset", err))
		}
	}

	return violations
}
