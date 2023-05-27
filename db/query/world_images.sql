
-- name: CreateWorldImages :exec
INSERT INTO world_images (world_id) VALUES (@world_id);

-- name: UpdateWorldImages :one
UPDATE world_images
SET
    image_header = COALESCE(sqlc.arg(image_header), image_header),
    image_avatar = COALESCE(sqlc.arg(image_avatar), image_avatar)
WHERE
    world_id = sqlc.arg(world_id)
RETURNING *;

-- name: DeleteWorldImages :exec
DELETE FROM world_images WHERE world_id = @world_id;

-- name: GetWorldImages :one
SELECT * FROM world_images WHERE world_id = @world_id LIMIT 1;