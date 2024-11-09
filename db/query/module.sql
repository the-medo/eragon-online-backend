-- name: GetModule :one
SELECT * FROM modules
WHERE
    (world_id IS NULL OR world_id = sqlc.narg(world_id)) AND
    (system_id IS NULL OR system_id = sqlc.narg(system_id)) AND
    (character_id IS NULL OR character_id = sqlc.narg(character_id)) AND
    (quest_id IS NULL OR quest_id = sqlc.narg(quest_id))
;

-- name: GetModuleById :one
SELECT * FROM view_modules WHERE id = sqlc.arg(module_id);

-- name: GetModulesByIDs :many
SELECT * FROM view_modules WHERE id = ANY(@module_ids::int[]);

-- name: CreateModule :one
INSERT INTO modules (module_type, menu_id, world_id, quest_id, character_id, system_id, description_post_id)
VALUES (sqlc.arg(module_type), sqlc.arg(menu_id), sqlc.narg(world_id), sqlc.narg(quest_id), sqlc.narg(character_id), sqlc.narg(system_id), sqlc.arg(description_post_id))
RETURNING *;

-- name: UpdateModule :one
UPDATE modules SET
    header_img_id = COALESCE(sqlc.narg(header_img_id), header_img_id),
    thumbnail_img_id = COALESCE(sqlc.narg(thumbnail_img_id), thumbnail_img_id),
    avatar_img_id = COALESCE(sqlc.narg(avatar_img_id), avatar_img_id),
    description_post_id = COALESCE(sqlc.narg(description_post_id), description_post_id)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteModule :exec
DELETE FROM modules WHERE id = sqlc.arg(id);