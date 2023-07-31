
-- name: CreateWorldMenu :one
INSERT INTO world_menu (world_id, menu_id)
VALUES (sqlc.arg(world_id), sqlc.arg(menu_id))
RETURNING *;

-- name: DeleteWorldMenu :exec
DELETE FROM world_menu WHERE world_id = sqlc.arg(world_id) AND menu_id = sqlc.arg(menu_id);

-- name: GetWorldMenu :one
SELECT * FROM world_menu WHERE world_id = sqlc.arg(world_id) AND menu_id = sqlc.arg(menu_id);