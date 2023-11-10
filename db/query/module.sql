-- name: GetModule :one
SELECT * FROM modules
WHERE
    world_id = COALESCE(sqlc.narg(world_id), world_id) OR
    quest_id = COALESCE(sqlc.narg(quest_id), quest_id) OR
    character_id = COALESCE(sqlc.narg(character_id), character_id) OR
    system_id = COALESCE(sqlc.narg(system_id), system_id)
;

-- name: GetModuleById :one
SELECT * FROM modules WHERE id = sqlc.arg(module_id);

-- name: GetModulesByIDs :many
SELECT * FROM modules WHERE id = ANY(@module_ids::int[]);

-- name: CreateModule :one
INSERT INTO modules (module_type, menu_id, world_id, quest_id, character_id, system_id)
VALUES (sqlc.arg(module_type), sqlc.arg(menu_id), sqlc.narg(world_id), sqlc.narg(quest_id), sqlc.narg(character_id), sqlc.narg(system_id))
RETURNING *;

-- name: UpdateModule :one
UPDATE modules SET
    header_img_id = COALESCE(sqlc.narg(header_img_id), header_img_id),
    thumbnail_img_id = COALESCE(sqlc.narg(thumbnail_img_id), thumbnail_img_id),
    avatar_img_id = COALESCE(sqlc.narg(avatar_img_id), avatar_img_id)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteModule :exec
DELETE FROM modules WHERE id = sqlc.arg(id);