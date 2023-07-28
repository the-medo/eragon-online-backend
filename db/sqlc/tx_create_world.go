package db

import (
	"context"
	"database/sql"
	"fmt"
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

		menu, err := q.CreateMenu(ctx, CreateMenuParams{
			MenuCode: "world-" + fmt.Sprint(world.ID),
			MenuHeaderImgID: sql.NullInt32{
				Int32: 0,
				Valid: false,
			},
		})
		if err != nil {
			return err
		}

		resourcesMenuItem, err := q.CreateMenuItem(ctx, CreateMenuItemParams{
			MenuID:            menu.ID,
			MenuItemCode:      "resources",
			Name:              "Resources",
			Position:          1,
			ParentItemID:      sql.NullInt32{},
			MenuItemImgID:     sql.NullInt32{},
			DescriptionPostID: sql.NullInt32{},
		})
		if err != nil {
			return err
		}

		_, err = q.CreateMenuItem(ctx, CreateMenuItemParams{
			MenuID:       menu.ID,
			MenuItemCode: "introduction",
			Name:         "Introduction",
			Position:     2,
			ParentItemID: sql.NullInt32{
				Int32: resourcesMenuItem.ID,
				Valid: true,
			},
			MenuItemImgID:     sql.NullInt32{},
			DescriptionPostID: sql.NullInt32{},
		})
		if err != nil {
			return err
		}

		_, err = q.CreateWorldMenu(ctx, CreateWorldMenuParams{
			WorldID: world.ID,
			MenuID:  menu.ID,
		})
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
			WorldID:            world.ID,
			UserID:             arg.UserId,
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
