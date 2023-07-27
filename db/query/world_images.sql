
-- name: CreateWorldImages :exec
INSERT INTO world_images (world_id) VALUES (@world_id);

-- name: UpdateWorldImages :one
UPDATE world_images
SET
    header_img_id = COALESCE(sqlc.arg(header_img_id), header_img_id),
    thumbnail_img_id = COALESCE(sqlc.arg(thumbnail_img_id), thumbnail_img_id),
    avatar_img_id = COALESCE(sqlc.arg(avatar_img_id), avatar_img_id)
WHERE
    world_id = sqlc.arg(world_id)
RETURNING *;

-- name: DeleteWorldImages :exec
DELETE FROM world_images WHERE world_id = @world_id;

-- name: GetWorldImages :one
SELECT * FROM world_images WHERE world_id = @world_id LIMIT 1;