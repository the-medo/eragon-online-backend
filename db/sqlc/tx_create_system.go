package db

import (
	"context"
	"database/sql"
	"fmt"
)

type CreateSystemTxParams struct {
	UserId int32
	CreateSystemParams
}

type CreateSystemTxResult struct {
	System *System
	Module *ViewModule
}

func (store *SQLStore) CreateSystemTx(ctx context.Context, arg CreateSystemTxParams) (CreateSystemTxResult, error) {
	var result CreateSystemTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		system, err := q.CreateSystem(ctx, arg.CreateSystemParams)
		if err != nil {
			return err
		}

		menu, err := q.CreateMenu(ctx, CreateMenuParams{
			MenuCode: "system-" + fmt.Sprint(system.ID),
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
			ModuleType: ModuleTypeSystem,
			MenuID:     menu.ID,
			SystemID: sql.NullInt32{
				Int32: system.ID,
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

		_, err = q.CreateModuleAdmin(ctx, CreateModuleAdminParams{
			ModuleID:           module.ID,
			UserID:             arg.UserId,
			SuperAdmin:         true,
			Approved:           1,
			MotivationalLetter: "Creator of the system!",
		})

		if err != nil {
			return err
		}

		mapPinTypeGroup, err := q.CreateMapPinTypeGroup(ctx, "System - "+fmt.Sprint(system.Name))
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

		result = CreateSystemTxResult{
			System: &system,
			Module: &ViewModule{
				ID:         module.ID,
				ModuleType: ModuleTypeSystem,
				SystemID:   sql.NullInt32{Int32: system.ID, Valid: true},
				MenuID:     menu.ID,
			},
		}

		return nil
	})

	return result, err
}
