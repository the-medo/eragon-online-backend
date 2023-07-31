package api

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/api/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"time"
)

func (server *Server) GetWorldMonthlyActivity(ctx context.Context, req *pb.GetWorldMonthlyActivityRequest) (*pb.GetWorldMonthlyActivityResponse, error) {

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

	worldActivity, err := server.store.GetWorldMonthlyActivity(ctx, arg)
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
