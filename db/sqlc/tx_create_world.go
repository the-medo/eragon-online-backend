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
	World  *World
	Module *ViewModule
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

		post, err := q.CreatePost(ctx, CreatePostParams{
			UserID:    arg.UserId,
			Title:     fmt.Sprintf("%s - introduction", arg.Name),
			IsDraft:   false,
			IsPrivate: false,
			Content:   "",
		})
		if err != nil {
			return err
		}

		module, err := q.CreateModule(ctx, CreateModuleParams{
			ModuleType: ModuleTypeWorld,
			MenuID:     menu.ID,
			WorldID: sql.NullInt32{
				Int32: world.ID,
				Valid: true,
			},
			DescriptionPostID: post.ID,
		})
		if err != nil {
			return err
		}

		_, err = q.CreateEntity(ctx, CreateEntityParams{
			Type:     EntityTypePost,
			ModuleID: module.ID,
			PostID: sql.NullInt32{
				Int32: post.ID,
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
			IsDefault:         sql.NullBool{Bool: true, Valid: true},
		})
		if err != nil {
			return err
		}

		result = CreateWorldTxResult{
			World: &world,
			Module: &ViewModule{
				ID:         module.ID,
				ModuleType: ModuleTypeWorld,
				WorldID:    sql.NullInt32{Int32: world.ID, Valid: true},
				MenuID:     menu.ID,
			},
		}

		return nil
	})

	return result, err
}
