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

func insertSubMenuItem(ctx context.Context, q *Queries, menuId int32, position int32, code string, name string, parentItemId int32) error {

	_, err := q.CreateMenuItem(ctx, CreateMenuItemParams{
		MenuID:       menuId,
		MenuItemCode: code,
		Name:         name,
		Position:     position,
		ParentItemID: sql.NullInt32{
			Int32: parentItemId,
			Valid: true,
		},
		DescriptionPostID: sql.NullInt32{},
	})
	if err != nil {
		return err
	}
	return nil
}

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

		overviewMenuItem, err := q.CreateMenuItem(ctx, CreateMenuItemParams{
			MenuID:            menu.ID,
			MenuItemCode:      "overview",
			Name:              "Overview",
			Position:          1,
			ParentItemID:      sql.NullInt32{},
			DescriptionPostID: sql.NullInt32{},
		})
		if err != nil {
			return err
		}

		err = insertSubMenuItem(ctx, q, menu.ID, 2, "races", "Races", overviewMenuItem.ID)
		if err != nil {
			return err
		}

		err = insertSubMenuItem(ctx, q, menu.ID, 3, "flora-and-fauna", "Flora & Fauna", overviewMenuItem.ID)
		if err != nil {
			return err
		}

		err = insertSubMenuItem(ctx, q, menu.ID, 4, "magic", "Magic", overviewMenuItem.ID)
		if err != nil {
			return err
		}

		err = insertSubMenuItem(ctx, q, menu.ID, 5, "science-and-technology", "Science & Technology", overviewMenuItem.ID)
		if err != nil {
			return err
		}

		err = insertSubMenuItem(ctx, q, menu.ID, 6, "history", "History", overviewMenuItem.ID)
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
			ID:               world.ID,
			Name:             world.Name,
			CreatedAt:        world.CreatedAt,
			Public:           world.Public,
			BasedOn:          world.BasedOn,
			ShortDescription: world.ShortDescription,
		}
		return nil
	})

	return result, err
}
