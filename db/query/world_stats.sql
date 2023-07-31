
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

-- name: GetWorldDailyActivity :many
SELECT
    *
FROM
    world_activity
WHERE
    world_id = COALESCE(sqlc.narg(world_id), world_id) AND
    date >= sqlc.arg(date_from)
ORDER BY date DESC;

-- name: GetWorldMonthlyActivity :many
SELECT
    cast(date_trunc('month', date) as date) AS month,
    world_id,
    cast(SUM(post_count) as integer) AS post_count,
    cast(SUM(quest_count) as integer) AS quest_count,
    cast(SUM(resource_count) as integer) AS resource_count
FROM
    world_activity
WHERE
    world_id = COALESCE(sqlc.narg(world_id), world_id) AND
    date >= sqlc.arg(date_from)
GROUP BY
    month, world_id
ORDER BY
    month DESC;