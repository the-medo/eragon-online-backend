package servicecore

import (
	"context"
	"database/sql"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (core *ServiceCore) SharedCreateMenuItem(ctx context.Context, req *pb.CreateMenuItemRequest) (*db.MenuItem, error) {
	argCreateEntityGroup := db.CreateEntityGroupParams{
		Name: sql.NullString{
			String: req.GetName(),
			Valid:  true,
		},
		Description: sql.NullString{},
		Style: sql.NullString{
			String: "framed",
			Valid:  true,
		},
		Direction: sql.NullString{
			String: "vertical",
			Valid:  true,
		},
	}

	newEntityGroup, err := core.Store.CreateEntityGroup(ctx, argCreateEntityGroup)
	if err != nil {
		return nil, err
	}

	arg := db.CreateMenuItemParams{
		MenuID:       req.GetMenuId(),
		MenuItemCode: req.GetCode(),
		Name:         req.GetName(),
		Position:     req.GetPosition(),
		IsMain: sql.NullBool{
			Bool:  req.GetIsMain(),
			Valid: req.IsMain != nil,
		},
		EntityGroupID: sql.NullInt32{
			Int32: newEntityGroup.ID,
			Valid: true,
		},
	}

	menuItem, err := core.Store.CreateMenuItem(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create menu item: %s", err)
	}

	return &menuItem, nil
}
