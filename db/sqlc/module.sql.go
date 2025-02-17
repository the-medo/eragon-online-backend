// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: module.sql

package db

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

const createModule = `-- name: CreateModule :one
INSERT INTO modules (module_type, menu_id, world_id, quest_id, character_id, system_id, description_post_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, module_type, menu_id, header_img_id, thumbnail_img_id, avatar_img_id, world_id, system_id, character_id, quest_id, description_post_id
`

type CreateModuleParams struct {
	ModuleType        ModuleType    `json:"module_type"`
	MenuID            int32         `json:"menu_id"`
	WorldID           sql.NullInt32 `json:"world_id"`
	QuestID           sql.NullInt32 `json:"quest_id"`
	CharacterID       sql.NullInt32 `json:"character_id"`
	SystemID          sql.NullInt32 `json:"system_id"`
	DescriptionPostID int32         `json:"description_post_id"`
}

func (q *Queries) CreateModule(ctx context.Context, arg CreateModuleParams) (Module, error) {
	row := q.db.QueryRowContext(ctx, createModule,
		arg.ModuleType,
		arg.MenuID,
		arg.WorldID,
		arg.QuestID,
		arg.CharacterID,
		arg.SystemID,
		arg.DescriptionPostID,
	)
	var i Module
	err := row.Scan(
		&i.ID,
		&i.ModuleType,
		&i.MenuID,
		&i.HeaderImgID,
		&i.ThumbnailImgID,
		&i.AvatarImgID,
		&i.WorldID,
		&i.SystemID,
		&i.CharacterID,
		&i.QuestID,
		&i.DescriptionPostID,
	)
	return i, err
}

const deleteModule = `-- name: DeleteModule :exec
DELETE FROM modules WHERE id = $1
`

func (q *Queries) DeleteModule(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteModule, id)
	return err
}

const getModule = `-- name: GetModule :one
SELECT id, module_type, menu_id, header_img_id, thumbnail_img_id, avatar_img_id, world_id, system_id, character_id, quest_id, description_post_id FROM modules
WHERE
    (world_id IS NULL OR world_id = $1) AND
    (system_id IS NULL OR system_id = $2) AND
    (character_id IS NULL OR character_id = $3) AND
    (quest_id IS NULL OR quest_id = $4)
`

type GetModuleParams struct {
	WorldID     sql.NullInt32 `json:"world_id"`
	SystemID    sql.NullInt32 `json:"system_id"`
	CharacterID sql.NullInt32 `json:"character_id"`
	QuestID     sql.NullInt32 `json:"quest_id"`
}

func (q *Queries) GetModule(ctx context.Context, arg GetModuleParams) (Module, error) {
	row := q.db.QueryRowContext(ctx, getModule,
		arg.WorldID,
		arg.SystemID,
		arg.CharacterID,
		arg.QuestID,
	)
	var i Module
	err := row.Scan(
		&i.ID,
		&i.ModuleType,
		&i.MenuID,
		&i.HeaderImgID,
		&i.ThumbnailImgID,
		&i.AvatarImgID,
		&i.WorldID,
		&i.SystemID,
		&i.CharacterID,
		&i.QuestID,
		&i.DescriptionPostID,
	)
	return i, err
}

const getModuleById = `-- name: GetModuleById :one
SELECT id, world_id, system_id, character_id, quest_id, module_type, menu_id, header_img_id, thumbnail_img_id, avatar_img_id, description_post_id, tags FROM view_modules WHERE id = $1
`

func (q *Queries) GetModuleById(ctx context.Context, moduleID int32) (ViewModule, error) {
	row := q.db.QueryRowContext(ctx, getModuleById, moduleID)
	var i ViewModule
	err := row.Scan(
		&i.ID,
		&i.WorldID,
		&i.SystemID,
		&i.CharacterID,
		&i.QuestID,
		&i.ModuleType,
		&i.MenuID,
		&i.HeaderImgID,
		&i.ThumbnailImgID,
		&i.AvatarImgID,
		&i.DescriptionPostID,
		pq.Array(&i.Tags),
	)
	return i, err
}

const getModulesByIDs = `-- name: GetModulesByIDs :many
SELECT id, world_id, system_id, character_id, quest_id, module_type, menu_id, header_img_id, thumbnail_img_id, avatar_img_id, description_post_id, tags FROM view_modules WHERE id = ANY($1::int[])
`

func (q *Queries) GetModulesByIDs(ctx context.Context, moduleIds []int32) ([]ViewModule, error) {
	rows, err := q.db.QueryContext(ctx, getModulesByIDs, pq.Array(moduleIds))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ViewModule{}
	for rows.Next() {
		var i ViewModule
		if err := rows.Scan(
			&i.ID,
			&i.WorldID,
			&i.SystemID,
			&i.CharacterID,
			&i.QuestID,
			&i.ModuleType,
			&i.MenuID,
			&i.HeaderImgID,
			&i.ThumbnailImgID,
			&i.AvatarImgID,
			&i.DescriptionPostID,
			pq.Array(&i.Tags),
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateModule = `-- name: UpdateModule :one
UPDATE modules SET
    header_img_id = COALESCE($1, header_img_id),
    thumbnail_img_id = COALESCE($2, thumbnail_img_id),
    avatar_img_id = COALESCE($3, avatar_img_id),
    description_post_id = COALESCE($4, description_post_id)
WHERE id = $5
RETURNING id, module_type, menu_id, header_img_id, thumbnail_img_id, avatar_img_id, world_id, system_id, character_id, quest_id, description_post_id
`

type UpdateModuleParams struct {
	HeaderImgID       sql.NullInt32 `json:"header_img_id"`
	ThumbnailImgID    sql.NullInt32 `json:"thumbnail_img_id"`
	AvatarImgID       sql.NullInt32 `json:"avatar_img_id"`
	DescriptionPostID sql.NullInt32 `json:"description_post_id"`
	ID                int32         `json:"id"`
}

func (q *Queries) UpdateModule(ctx context.Context, arg UpdateModuleParams) (Module, error) {
	row := q.db.QueryRowContext(ctx, updateModule,
		arg.HeaderImgID,
		arg.ThumbnailImgID,
		arg.AvatarImgID,
		arg.DescriptionPostID,
		arg.ID,
	)
	var i Module
	err := row.Scan(
		&i.ID,
		&i.ModuleType,
		&i.MenuID,
		&i.HeaderImgID,
		&i.ThumbnailImgID,
		&i.AvatarImgID,
		&i.WorldID,
		&i.SystemID,
		&i.CharacterID,
		&i.QuestID,
		&i.DescriptionPostID,
	)
	return i, err
}
