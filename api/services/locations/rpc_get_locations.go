package locations

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

func (server *ServiceLocations) GetLocations(ctx context.Context, req *pb.GetLocationsRequest) (*pb.GetLocationsResponse, error) {
	violations := validateGetLocations(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	limit, offset := apihelpers.GetDefaultQueryBoundaries(req.GetLimit(), req.GetOffset())

	arg := db.GetLocationsParams{
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

	locationRows, err := server.Store.GetLocations(ctx, arg)

	if err != nil {
		return nil, err
	}

	rsp := &pb.GetLocationsResponse{
		Locations:  make([]*pb.Location, len(locationRows)),
		TotalCount: 0,
	}

	if len(locationRows) > 0 {
		rsp.TotalCount = locationRows[0].TotalCount
	}

	for i, locationRow := range locationRows {
		rsp.Locations[i] = converters.ConvertGetLocationsRow(locationRow)
	}

	return rsp, nil
}

func validateGetLocations(req *pb.GetLocationsRequest) (violations []*errdetails.BadRequest_FieldViolation) {
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
