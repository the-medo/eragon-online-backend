-- name: GetWorldTagsAvailable :many
SELECT * FROM world_tags_available;

-- name: GetWorldTagAvailable :one
SELECT * FROM world_tags_available WHERE id = sqlc.arg(tag_id);

-- name: CreateWorldTagAvailable :one
INSERT INTO world_tags_available (tag)
VALUES (sqlc.arg(tag))
RETURNING *;

-- name: UpdateWorldTagAvailable :one
UPDATE world_tags_available
SET tag = sqlc.arg(tag)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteWorldTagAvailable :exec
DELETE FROM world_tags_available WHERE id = sqlc.arg(id);

-- name: GetWorldTags :many
SELECT * FROM world_tags;

-- name: GetWorldTag :one
SELECT * FROM world_tags WHERE world_id = sqlc.arg(world_id) AND tag_id = sqlc.arg(tag_id);

-- name: CreateWorldTag :one
INSERT INTO world_tags (world_id, tag_id)
VALUES (sqlc.arg(world_id), sqlc.arg(tag_id))
RETURNING *;

-- name: DeleteWorldTag :exec
DELETE FROM world_tags WHERE world_id = sqlc.arg(world_id) AND tag_id = sqlc.arg(tag_id);
