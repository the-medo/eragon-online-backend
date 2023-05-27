-- name: CreateWorld :one
INSERT INTO worlds (
    name,
    description
) VALUES (
     @name, @description
 ) RETURNING *;

-- name: UpdateWorld :one
UPDATE worlds
SET
    name = COALESCE(sqlc.arg(name), name),
    public = COALESCE(sqlc.arg(public), public),
    description = COALESCE(sqlc.arg(description), description)
WHERE
    id = sqlc.arg(world_id)
RETURNING *;

-- name: DeleteWorld :exec
DELETE FROM worlds WHERE id = @world_id;

-- name: GetWorldByID :one
SELECT * FROM view_worlds WHERE id = @world_id LIMIT 1;

-- name: GetWorlds :many
SELECT * FROM view_worlds
WHERE (@is_public::boolean IS NULL OR public = @is_public)
ORDER BY
    CASE
     WHEN @order_result::bool
         THEN @order_by::VARCHAR
     ELSE 'activity'
     END
DESC
LIMIT @page_limit
OFFSET @page_offset;