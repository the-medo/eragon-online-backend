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

-- name: GetModuleEntityTagsAvailable :many
SELECT * FROM module_entity_tags_available
WHERE module_id = sqlc.arg(module_id);

-- name: CreateModuleEntityTagAvailable :one
INSERT INTO module_entity_tags_available (module_id, tag)
VALUES (sqlc.arg(module_id), sqlc.arg(tag))
RETURNING *;

-- name: UpdateModuleEntityTagAvailable :one
UPDATE module_entity_tags_available
SET tag = sqlc.arg(tag)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteModuleEntityTagAvailable :exec
DELETE FROM module_entity_tags_available WHERE id = sqlc.arg(id);

-- name: CreateEntityTag :one
INSERT INTO entity_tags (entity_id, tag_id)
VALUES (sqlc.arg(entity_id), sqlc.arg(tag_id))
RETURNING *;

-- name: DeleteEntityTag :exec
DELETE FROM entity_tags
WHERE
        entity_id = COALESCE(sqlc.narg(entity_id), entity_id) AND
        tag_id = COALESCE(sqlc.narg(tag_id), tag_id) AND
    (NOT (sqlc.narg(entity_id) IS NULL AND sqlc.narg(tag_id) IS NULL));