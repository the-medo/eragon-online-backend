-- name: CreateLocation :one
INSERT INTO locations (name, description, post_id, thumbnail_image_id)
VALUES (sqlc.arg(name), sqlc.narg(description), sqlc.narg(post_id), sqlc.narg(thumbnail_image_id))
RETURNING *;

-- name: GetLocations :many
SELECT * FROM view_locations;

-- name: GetLocationsByModule :many
SELECT * FROM view_locations
WHERE module_id = sqlc.arg(module_id)
;

-- name: GetLocationsByIDs :many
SELECT * FROM locations WHERE id = ANY(@location_ids::int[]);

-- name: GetLocationById :one
SELECT * FROM locations WHERE id = sqlc.arg(id);

-- name: GetViewLocationById :one
SELECT * FROM view_locations WHERE id = sqlc.arg(id);

-- name: UpdateLocation :one
UPDATE locations
SET
    name = COALESCE(sqlc.narg(name), name),
    description = COALESCE(sqlc.narg(description), description),
    post_id = COALESCE(sqlc.narg(post_id), post_id),
    thumbnail_image_id = COALESCE(sqlc.narg(thumbnail_image_id), thumbnail_image_id)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteLocation :exec
CALL delete_location(sqlc.arg(id));
