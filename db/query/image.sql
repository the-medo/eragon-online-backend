-- name: CreateImage :one
INSERT INTO images
(
    img_guid,
    image_type_id,
    name,
    url,
    base_url,
    user_id,
    width,
    height
)
VALUES
    (@img_guid, @img_type_id, @name, @url, @base_url, @user_id, @width, @height)
RETURNING *;

-- name: GetImageById :one
SELECT * FROM images WHERE id = @id LIMIT 1;

-- name: GetImagesByIDs :many
SELECT * FROM images WHERE id = ANY(@image_ids::int[]);

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
    base_url = COALESCE(sqlc.arg(base_url), base_url),
    user_id = COALESCE(sqlc.arg(user_id), user_id),
    width = COALESCE(sqlc.arg(width), width),
    height = COALESCE(sqlc.arg(height), height)
WHERE
    id = sqlc.arg(id)
RETURNING *;

-- name: DeleteImage :exec
DELETE FROM images WHERE id = @id;

-- name: GetImages :many
SELECT
    *
FROM
    images
WHERE
    (user_id = COALESCE(sqlc.narg(user_id), user_id)) AND
    (image_type_id = COALESCE(sqlc.narg(image_type_id), image_type_id))
ORDER BY id DESC
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: GetImagesCount :one
SELECT COUNT(*) FROM images
WHERE
    (user_id = COALESCE(sqlc.narg(user_id), user_id)) AND
    (image_type_id = COALESCE(sqlc.narg(image_type_id), image_type_id));