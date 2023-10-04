-- name: GetEntityIDsOfGroup :one
WITH entity_data AS (
    SELECT
        e.post_id,
        e.map_id,
        e.location_id,
        e.image_id
    FROM
        get_recursive_entities(sqlc.arg(entity_group_id)) er
        JOIN entities e ON er.content_entity_id = e.id
)
SELECT
    CAST(ARRAY_AGG(post_id) FILTER (WHERE post_id IS NOT NULL) as INT[]) AS post_ids,
    CAST(ARRAY_AGG(map_id) FILTER (WHERE map_id IS NOT NULL) as INT[]) AS map_ids,
    CAST(ARRAY_AGG(location_id) FILTER (WHERE location_id IS NOT NULL) as INT[]) AS location_ids,
    CAST(ARRAY_AGG(image_id) FILTER (WHERE image_id IS NOT NULL) as INT[]) AS image_ids
FROM
    entity_data;

-- name: GetEntityGroupContents :many
WITH entity_data AS (
    SELECT
        eg.id as entity_group_id,
        eg.name as entity_group_name,
        eg.description as entity_group_description,
        e.post_id,
        e.map_id,
        e.location_id,
        e.image_id
    FROM
        get_recursive_entities(sqlc.arg(entity_group_id)) er
            JOIN entities e ON er.content_entity_id = e.id
            JOIN entity_groups eg ON er.entity_group_id = eg.id
)
SELECT
    entity_group_id,
    CAST(min(entity_group_name) as VARCHAR) as entity_group_name,
    CAST(min(entity_group_description) as VARCHAR) as entity_group_description,
    CAST(ARRAY_AGG(post_id) FILTER (WHERE post_id IS NOT NULL) as INT[]) AS post_ids,
    CAST(ARRAY_AGG(map_id) FILTER (WHERE map_id IS NOT NULL) as INT[]) AS map_ids,
    CAST(ARRAY_AGG(location_id) FILTER (WHERE location_id IS NOT NULL) as INT[]) AS location_ids,
    CAST(ARRAY_AGG(image_id) FILTER (WHERE image_id IS NOT NULL) as INT[]) AS image_ids
FROM
    entity_data
GROUP BY
    entity_group_id;

-- name: CreateEntity :one
INSERT INTO entities (type, post_id, map_id, location_id, image_id)
VALUES (sqlc.arg(type), sqlc.narg(post_id), sqlc.narg(map_id), sqlc.narg(location_id), sqlc.narg(image_id))
RETURNING *;

-- name: GetEntityByID :one
SELECT * FROM entities WHERE id = sqlc.arg(id);

-- name: UpdateEntity :one
UPDATE entities
SET
    type = COALESCE(sqlc.narg(type), type),
    post_id = COALESCE(sqlc.narg(post_id), post_id),
    map_id = COALESCE(sqlc.narg(map_id), map_id),
    location_id = COALESCE(sqlc.narg(location_id), location_id),
    image_id = COALESCE(sqlc.narg(image_id), image_id)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteEntity :exec
DELETE FROM entities WHERE id = sqlc.arg(id);

-- name: CreateEntityGroup :one
INSERT INTO entity_groups (name, description)
VALUES (sqlc.arg(name), sqlc.arg(description))
RETURNING *;

-- name: GetEntityGroupByID :one
SELECT * FROM entity_groups WHERE id = sqlc.arg(id);

-- name: UpdateEntityGroup :one
UPDATE entity_groups
SET
    name = COALESCE(sqlc.narg(name), name),
    description = COALESCE(sqlc.narg(description), description)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteEntityGroup :exec
DELETE FROM entity_groups WHERE id = sqlc.arg(id);

-- name: CreateEntityGroupContent :one
INSERT INTO entity_group_content (entity_group_id, position, content_entity_id, content_entity_group_id)
VALUES (sqlc.arg(entity_group_id), sqlc.arg(position), sqlc.narg(content_entity_id), sqlc.narg(content_entity_group_id))
RETURNING *;

-- name: GetEntityGroupContentByID :one
SELECT * FROM entity_group_content WHERE id = sqlc.arg(id);

-- name: UpdateEntityGroupContent :one
UPDATE entity_group_content
SET
    entity_group_id = COALESCE(sqlc.narg(entity_group_id), entity_group_id),
    position = COALESCE(sqlc.narg(position), position),
    content_entity_id = COALESCE(sqlc.narg(content_entity_id), content_entity_id),
    content_entity_group_id = COALESCE(sqlc.narg(content_entity_group_id), content_entity_group_id)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteEntityGroupContent :exec
DELETE FROM entity_group_content WHERE id = sqlc.arg(id);