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

func (server *Server) GetWorldDailyActivity(ctx context.Context, req *pb.GetWorldDailyActivityRequest) (*pb.GetWorldDailyActivityResponse, error) {

	violations := validateGetWorldDailyActivity(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	arg := db.GetWorldDailyActivityParams{
		WorldID: sql.NullInt32{
			Int32: req.GetWorldId(),
			Valid: req.WorldId != nil,
		},
		DateFrom: time.Now().AddDate(0, -1, 0),
	}

	if req.DateFrom != nil {
		arg.DateFrom = req.GetDateFrom().AsTime()
	}

	worldActivity, err := server.Store.GetWorldDailyActivity(ctx, arg)
	if err != nil {
		return nil, err
	}

	rsp := &pb.GetWorldDailyActivityResponse{
		Activity: make([]*pb.WorldActivity, len(worldActivity)),
	}

	for i, activity := range worldActivity {
		rsp.Activity[i] = converters.ConvertWorldActivity(activity)
	}

	return rsp, nil
}

func validateGetWorldDailyActivity(req *pb.GetWorldDailyActivityRequest) (violations []*errdetails.BadRequest_FieldViolation) {
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
