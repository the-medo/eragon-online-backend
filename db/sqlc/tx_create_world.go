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

		err = q.CreateWorldActivity(ctx, CreateWorldActivityParams{
			WorldID: world.ID,
			Date:    world.CreatedAt,
		})
		if err != nil {
			return err
		}

		err = q.CreateWorldImages(ctx, world.ID)
		if err != nil {
			return err
		}

		_, err = q.InsertWorldAdmin(ctx, InsertWorldAdminParams{
			WorldID: sql.NullInt32{
				Int32: world.ID,
				Valid: true,
			},
			UserID: sql.NullInt32{
				Int32: arg.UserId,
				Valid: true,
			},
			SuperAdmin:         true,
			Approved:           1,
			MotivationalLetter: "Creator of the world!",
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
			BasedOn:     world.BasedOn,
		}
		return nil
	})

	return result, err
}
