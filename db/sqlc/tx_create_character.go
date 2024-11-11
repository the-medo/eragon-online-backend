package db

import (
	"context"
	"database/sql"
	"fmt"
)

type CreateCharacterTxParams struct {
	UserId int32
	CreateCharacterParams
}

type CreateCharacterTxResult struct {
	Character *Character
	Module    *ViewModule
}

func (store *SQLStore) CreateCharacterTx(ctx context.Context, arg CreateCharacterTxParams) (CreateCharacterTxResult, error) {
	var result CreateCharacterTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		character, err := q.CreateCharacter(ctx, arg.CreateCharacterParams)
		if err != nil {
			return err
		}

		menu, err := q.CreateMenu(ctx, CreateMenuParams{
			MenuCode: "character-" + fmt.Sprint(character.ID),
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
			ModuleType: ModuleTypeCharacter,
			MenuID:     menu.ID,
			CharacterID: sql.NullInt32{
				Int32: character.ID,
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
			MotivationalLetter: "Creator of the character!",
		})

		if err != nil {
			return err
		}

		mapPinTypeGroup, err := q.CreateMapPinTypeGroup(ctx, "Character - "+fmt.Sprint(character.Name))
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

		result = CreateCharacterTxResult{
			Character: &character,
			Module: &ViewModule{
				ID:          module.ID,
				ModuleType:  ModuleTypeCharacter,
				CharacterID: sql.NullInt32{Int32: character.ID, Valid: true},
				MenuID:      menu.ID,
			},
		}

		return nil
	})

	return result, err
}
