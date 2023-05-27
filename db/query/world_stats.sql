
-- name: CreateWorldStats :exec
INSERT INTO world_stats ( world_id ) VALUES ( @world_id );

-- name: UpdateWorldStats :one
UPDATE world_stats
SET
    final_content_rating = COALESCE(sqlc.arg(final_content_rating), final_content_rating),
    final_activity = COALESCE(sqlc.arg(final_activity), final_activity)
WHERE
        world_id = sqlc.arg(world_id)
RETURNING *;

-- name: DeleteWorldStats :exec
DELETE FROM world_stats
WHERE world_id = $1;

-- name: GetWorldStats :one
SELECT * FROM world_stats WHERE world_id = @world_id LIMIT 1;

-- name: InsertWorldStatsHistory :one
INSERT INTO world_stats_history
(
    world_id,
    final_content_rating,
    final_activity
)
VALUES
    (@world_id, @final_content_rating, @final_activity)
RETURNING *;

-- name: DeleteWorldStatsHistory :exec
DELETE FROM world_stats_history WHERE world_id = @world_id;

-- name: GetWorldStatsHistory :many
SELECT * FROM world_stats_history
WHERE created_at >= @start_date
ORDER BY created_at DESC;