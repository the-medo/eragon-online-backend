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

type CreateWorldTxResult struct {
	ViewWorld *ViewWorld
	Module    *ViewModule
}

func insertSubMenuItem(ctx context.Context, q *Queries, menuId int32, position int32, code string, name string, isMain bool) error {

	_, err := q.CreateMenuItem(ctx, CreateMenuItemParams{
		MenuID:       menuId,
		MenuItemCode: code,
		Name:         name,
		Position:     position,
		IsMain: sql.NullBool{
			Bool:  isMain,
			Valid: true,
		},
		DescriptionPostID: sql.NullInt32{},
	})
	if err != nil {
		return err
	}
	return nil
}

func (store *SQLStore) CreateWorldTx(ctx context.Context, arg CreateWorldTxParams) (CreateWorldTxResult, error) {
	var result CreateWorldTxResult

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

		module, err := q.CreateModule(ctx, CreateModuleParams{
			ModuleType: ModuleTypeWorld,
			MenuID: sql.NullInt32{
				Int32: menu.ID,
				Valid: true,
			},
			WorldID: sql.NullInt32{
				Int32: world.ID,
				Valid: true,
			},
		})
		if err != nil {
			return err
		}

		_, err = q.UpsertUserModule(ctx, UpsertUserModuleParams{
			ModuleID: module.ID,
			UserID:   arg.UserId,
			Admin: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
			Following: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
			Favorite: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
		})
		if err != nil {
			return err
		}

		_, err = q.CreateMenuItem(ctx, CreateMenuItemParams{
			MenuID:       menu.ID,
			MenuItemCode: "overview",
			Name:         "Overview",
			Position:     1,
			IsMain: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
			DescriptionPostID: sql.NullInt32{},
		})
		if err != nil {
			return err
		}

		err = insertSubMenuItem(ctx, q, menu.ID, 2, "races", "Races", false)
		if err != nil {
			return err
		}

		err = insertSubMenuItem(ctx, q, menu.ID, 3, "flora-and-fauna", "Flora & Fauna", false)
		if err != nil {
			return err
		}

		err = insertSubMenuItem(ctx, q, menu.ID, 4, "magic", "Magic", false)
		if err != nil {
			return err
		}

		err = insertSubMenuItem(ctx, q, menu.ID, 5, "science-and-technology", "Science & Technology", false)
		if err != nil {
			return err
		}

		err = insertSubMenuItem(ctx, q, menu.ID, 6, "history", "History", false)
		if err != nil {
			return err
		}

		_, err = q.CreateModuleAdmin(ctx, CreateModuleAdminParams{
			ModuleID:           module.ID,
			UserID:             arg.UserId,
			SuperAdmin:         true,
			Approved:           1,
			MotivationalLetter: "Creator of the world!",
		})

		if err != nil {
			return err
		}

		mapPinTypeGroup, err := q.CreateMapPinTypeGroup(ctx, "World - "+fmt.Sprint(world.Name))
		if err != nil {
			return err
		}

		_, err = q.CreateModuleMapPinTypeGroup(ctx, CreateModuleMapPinTypeGroupParams{
			ModuleID:          module.ID,
			MapPinTypeGroupID: mapPinTypeGroup.ID,
		})
		if err != nil {
			return err
		}

		_, err = q.CreateMapPinType(ctx, CreateMapPinTypeParams{
			MapPinTypeGroupID: mapPinTypeGroup.ID,
			Shape:             PinShapePin,
			BackgroundColor:   sql.NullString{String: "#ffffff", Valid: true},
			BorderColor:       sql.NullString{String: "#000000", Valid: true},
			IconColor:         sql.NullString{String: "#000000", Valid: true},
			Width:             sql.NullInt32{Int32: 24, Valid: true},
			Section:           "Base",
		})
		if err != nil {
			return err
		}

		result = CreateWorldTxResult{
			ViewWorld: &ViewWorld{
				ID:               world.ID,
				Name:             world.Name,
				CreatedAt:        world.CreatedAt,
				Public:           world.Public,
				BasedOn:          world.BasedOn,
				ShortDescription: world.ShortDescription,
				ModuleID:         module.ID,
				MenuID:           sql.NullInt32{Int32: menu.ID, Valid: true},
			},
			Module: &ViewModule{
				ModuleID:      module.ID,
				ModuleType:    ModuleTypeWorld,
				ModuleWorldID: sql.NullInt32{Int32: world.ID, Valid: true},
				MenuID:        sql.NullInt32{Int32: menu.ID, Valid: true},
			},
		}

		return nil
	})

	return result, err
}
