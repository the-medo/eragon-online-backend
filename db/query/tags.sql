-- name: GetModuleTypeTagsAvailable :many
SELECT * FROM view_module_type_tags_available
WHERE module_type = sqlc.arg(module_type);

-- name: GetModuleTypeTagAvailable :one
SELECT * FROM view_module_type_tags_available WHERE id = sqlc.arg(tag_id);

-- name: CreateModuleTypeTagAvailable :one
INSERT INTO module_type_tags_available (module_type, tag)
VALUES (sqlc.arg(module_type), sqlc.arg(tag))
RETURNING *;

-- name: UpdateModuleTypeTagAvailable :one
UPDATE module_type_tags_available
SET tag = sqlc.arg(tag)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteModuleTypeTagAvailable :exec
DELETE FROM module_type_tags_available WHERE id = sqlc.arg(id);

-- name: CreateModuleTag :one
INSERT INTO module_tags (module_id, tag_id)
VALUES (sqlc.arg(module_id), sqlc.arg(tag_id))
RETURNING *;

-- name: DeleteModuleTag :exec
DELETE FROM module_tags
WHERE
        module_id = COALESCE(sqlc.narg(module_id), module_id) AND
        tag_id = COALESCE(sqlc.narg(tag_id), tag_id) AND
    (NOT (sqlc.narg(module_id) IS NULL AND sqlc.narg(tag_id) IS NULL));
