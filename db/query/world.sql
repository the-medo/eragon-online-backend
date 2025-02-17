-- name: CreateWorld :one
INSERT INTO worlds (
    name,
    based_on,
    short_description
) VALUES (
     @name, @based_on, @short_description
 ) RETURNING *;

-- name: UpdateWorld :one
UPDATE worlds
SET
    name = COALESCE(sqlc.narg(name), name),
    based_on = COALESCE(sqlc.narg(based_on), based_on),
    public = COALESCE(sqlc.narg(public), public),
    short_description = COALESCE(sqlc.narg(short_description), short_description)
WHERE
    id = sqlc.arg(world_id)
RETURNING *;

-- name: DeleteWorld :exec
DELETE FROM worlds WHERE id = @world_id;

-- name: GetWorldByID :one
SELECT * FROM worlds WHERE id = @world_id LIMIT 1;

-- name: GetWorldsByIDs :many
SELECT * FROM worlds WHERE id = ANY(@world_ids::int[]);

-- name: GetWorlds :many
SELECT * FROM get_worlds(@is_public::boolean, @tags::integer[], @order_by::VARCHAR, 'DESC', @page_limit, @page_offset);

-- name: GetWorldsCount :one
SELECT COUNT(*) FROM view_worlds
WHERE (@is_public::boolean IS NULL OR public = @is_public) AND
    (array_length(@tags::integer[], 1) IS NULL OR tags @> @tags::integer[]);