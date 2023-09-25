package api

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/api/e"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"time"
)

func (server *Server) GetWorldMonthlyActivity(ctx context.Context, req *pb.GetWorldMonthlyActivityRequest) (*pb.GetWorldMonthlyActivityResponse, error) {

	violations := validateGetWorldMonthlyActivity(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	arg := db.GetWorldMonthlyActivityParams{
		WorldID: sql.NullInt32{
			Int32: req.GetWorldId(),
			Valid: req.WorldId != nil,
		},
		DateFrom: time.Now().AddDate(-1, 0, 0),
	}

	if req.DateFrom != nil {
		arg.DateFrom = req.GetDateFrom().AsTime()
	}

	worldActivity, err := server.Store.GetWorldMonthlyActivity(ctx, arg)
	if err != nil {
		return nil, err
	}

	rsp := &pb.GetWorldMonthlyActivityResponse{
		Activity: make([]*pb.WorldActivity, len(worldActivity)),
	}

	for i, activity := range worldActivity {
		rsp.Activity[i] = converters.ConvertMonthlyWorldActivity(activity)
	}

	return rsp, nil
}

func validateGetWorldMonthlyActivity(req *pb.GetWorldMonthlyActivityRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if req.WorldId != nil {
		if err := validator.ValidateWorldId(req.GetWorldId()); err != nil {
			violations = append(violations, e.FieldViolation("world_id", err))
		}
	}

	if req.DateFrom != nil {
		if err := validator.ValidateDatePast(req.GetDateFrom()); err != nil {
			violations = append(violations, e.FieldViolation("date_from", err))
		}

	}

	return violations
}
