-- name: CreateImage :one
INSERT INTO images
(
    img_guid,
    image_type_id,
    name,
    url,
    base_url
)
VALUES
    (@img_guid, @img_type_id, @name, @url, @base_url)
RETURNING *;

-- name: GetImageById :one
SELECT * FROM images WHERE id = @id LIMIT 1;

-- name: GetImageByGUID :one
SELECT * FROM images WHERE img_guid = @img_guid LIMIT 1;

-- name: GetImagesByImageTypeId :many
SELECT * FROM images WHERE image_type_id = @img_type_id;

-- name: UpdateImage :one
UPDATE images
SET
    img_guid = COALESCE(sqlc.arg(img_guid), img_guid),
    image_type_id = COALESCE(sqlc.arg(image_type_id), image_type_id),
    name = COALESCE(sqlc.arg(name), name),
    url = COALESCE(sqlc.arg(url), url),
    base_url = COALESCE(sqlc.arg(base_url), base_url)
WHERE
    id = sqlc.arg(id)
RETURNING *;

-- name: DeleteImage :exec
DELETE FROM images WHERE id = @id;