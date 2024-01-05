package maps

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/the-medo/talebound-backend/api/apihelpers"
	"github.com/the-medo/talebound-backend/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceMaps) GetMaps(ctx context.Context, req *pb.GetMapsRequest) (*pb.GetMapsResponse, error) {
	violations := validateGetMaps(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	limit, offset := apihelpers.GetDefaultQueryBoundaries(req.GetLimit(), req.GetOffset())

	arg := db.GetMapsParams{
		PageLimit:  limit,
		PageOffset: offset,
		ModuleID: sql.NullInt32{
			Int32: req.GetModuleId(),
			Valid: req.ModuleId != nil,
		},
		OrderBy: sql.NullString{
			String: req.GetOrderBy(),
			Valid:  req.OrderBy != nil,
		},
		Tags: req.GetTags(),
	}

	mapRows, err := server.Store.GetMaps(ctx, arg)

	if err != nil {
		return nil, err
	}

	rsp := &pb.GetMapsResponse{
		Maps:       make([]*pb.Map, len(mapRows)),
		TotalCount: 0,
	}

	if len(mapRows) > 0 {
		rsp.TotalCount = mapRows[0].TotalCount
	}

	for i, mapRow := range mapRows {
		rsp.Maps[i] = converters.ConvertGetMapsRow(mapRow)
	}

	return rsp, nil
}

func validateGetMaps(req *pb.GetMapsRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	fields := []string{"id", "name", "description"}

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

	if req.ModuleId != nil {
		if err := validator.ValidateUniversalId(req.GetModuleId()); err != nil {
			violations = append(violations, e.FieldViolation("module_id", err))
		}
	}

	return violations
}
