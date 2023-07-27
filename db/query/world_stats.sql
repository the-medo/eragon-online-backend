
-- name: CreateWorldActivity :exec
INSERT INTO world_activity ( world_id, date, post_count, quest_count, resource_count ) VALUES ( @world_id, @date, 0, 0, 0 );

-- name: UpdateWorldActivity :one
UPDATE world_activity
SET
    post_count = COALESCE(sqlc.narg(post_count), post_count),
    quest_count = COALESCE(sqlc.narg(quest_count), quest_count),
    resource_count = COALESCE(sqlc.narg(resource_count), resource_count)
WHERE
    world_id = sqlc.arg(world_id) AND
    date = sqlc.arg(date)
RETURNING *;

-- name: DeleteAllWorldActivity :exec
DELETE FROM world_activity WHERE world_id = @world_id;

-- name: DeleteWorldActivityForDate :exec
DELETE FROM world_activity WHERE world_id = @world_id AND date = @date;

-- name: GetWorldActivity :many
SELECT * FROM world_activity WHERE world_id = @world_id ORDER BY date DESC LIMIT 30;
