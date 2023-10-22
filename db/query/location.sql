-- name: CreateLocation :one
INSERT INTO locations (name, description, post_id, thumbnail_image_id)
VALUES (sqlc.arg(name), sqlc.narg(description), sqlc.narg(post_id), sqlc.narg(thumbnail_image_id))
RETURNING *;

-- name: GetLocations :many
SELECT * FROM view_locations;

-- name: GetLocationsForPlacement :many
SELECT
    vl.*
FROM
    view_locations vl
    LEFT JOIN world_locations wl ON vl.id = wl.location_id
    --LEFT JOIN quest_locations ql ON vl.id = ql.location_id
WHERE wl.world_id = sqlc.arg(world_id) --OR ql.quest_id = sqlc.arg(quest_id);
;

-- name: GetLocationByID :one
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

-- name: CreateWorldLocation :one
INSERT INTO world_locations (world_id, location_id)
VALUES (sqlc.arg(world_id), sqlc.arg(location_id))
RETURNING *;

-- name: GetWorldLocations :many
SELECT l.*
FROM view_locations l
    JOIN world_locations wl ON l.id = wl.location_id
WHERE wl.world_id = sqlc.arg(world_id);

-- name: GetLocationAssignments :one
SELECT
    CAST(MAX(COALESCE(wl.world_id, 0)) as integer) AS world_id,
    0 AS quest_id
FROM
    locations l
    LEFT JOIN world_locations wl ON l.id = wl.location_id
WHERE l.id = sqlc.arg(location_id)
GROUP BY l.id;

-- name: DeleteWorldLocation :exec
DELETE FROM world_locations
WHERE world_id = sqlc.arg(world_id) AND location_id = sqlc.arg(location_id);