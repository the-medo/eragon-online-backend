-- name: GetImageTypeById :one
SELECT * FROM image_types WHERE id = @id LIMIT 1;

-- name: GetImageTypeByName :one
SELECT * FROM image_types WHERE name = @name LIMIT 1;