package db

import (
	"context"
	"database/sql"
)

type CreateWorldTxParams struct {
	UserId int32
	CreateWorldParams
}

//
//// CreateUserTxResult is the result of the transfer transaction
//type CreateWorldTxResult struct {
//	world ViewWorld `json:"world"`
//}

func (store *SQLStore) CreateWorldTx(ctx context.Context, arg CreateWorldTxParams) (ViewWorld, error) {
	var result ViewWorld

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		world, err := q.CreateWorld(ctx, arg.CreateWorldParams)
		if err != nil {
			return err
		}

		err = q.CreateWorldActivity(ctx, world.ID)
		if err != nil {
			return err
		}

		err = q.CreateWorldImages(ctx, world.ID)
		if err != nil {
			return err
		}

		_, err = q.CreateWorldAdmin(ctx, CreateWorldAdminParams{
			WorldID: sql.NullInt32{
				Int32: world.ID,
				Valid: true,
			},
			UserID: sql.NullInt32{
				Int32: arg.UserId,
				Valid: true,
			},
			IsMain: true,
		})

		if err != nil {
			return err
		}

		result = ViewWorld{
			ID:          world.ID,
			Name:        world.Name,
			Description: world.Description,
			CreatedAt:   world.CreatedAt,
			Public:      world.Public,
		}
		return nil
	})

	return result, err
}
