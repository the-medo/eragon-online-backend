package api

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/api/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"time"
)

func (server *Server) GetWorldDailyActivity(ctx context.Context, req *pb.GetWorldDailyActivityRequest) (*pb.GetWorldDailyActivityResponse, error) {

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

	worldActivity, err := server.store.GetWorldDailyActivity(ctx, arg)
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
