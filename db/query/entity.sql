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