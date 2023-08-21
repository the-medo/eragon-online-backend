-- name: CreateMenu :one
INSERT INTO menus (menu_code, menu_header_img_id)
VALUES (sqlc.arg(menu_code), sqlc.narg(menu_header_img_id))
RETURNING *;

-- name: UpdateMenu :one
UPDATE menus
SET menu_code = COALESCE(sqlc.narg(menu_code), menu_code),
    menu_header_img_id = COALESCE(sqlc.narg(menu_header_img_id), menu_header_img_id)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteMenu :exec
DELETE FROM menus WHERE id = sqlc.arg(id);

-- name: GetMenu :one
SELECT * FROM menus WHERE id = sqlc.arg(id);

-- name: CreateMenuItem :one
INSERT INTO menu_items (menu_id, menu_item_code, name, position, is_main, description_post_id)
VALUES (sqlc.arg(menu_id), sqlc.arg(menu_item_code), sqlc.arg(name), sqlc.arg(position), sqlc.narg(is_main), sqlc.narg(description_post_id))
RETURNING *;

-- name: UpdateMenuItem :one
UPDATE menu_items
SET
    menu_item_code = COALESCE(sqlc.narg(menu_item_code), menu_item_code),
    name = COALESCE(sqlc.narg(name), name),
    position = COALESCE(sqlc.narg(position), position),
    is_main = COALESCE(sqlc.narg(is_main), is_main),
    description_post_id = COALESCE(sqlc.narg(description_post_id), description_post_id)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteMenuItem :exec
DELETE FROM menu_items WHERE id = sqlc.arg(id);

-- name: GetMenuItems :many
SELECT * FROM menu_items WHERE menu_id = sqlc.arg(menu_id);

-- name: MenuItemChangePositions :exec
UPDATE menu_items SET position = position + sqlc.arg(amount) WHERE menu_id = sqlc.arg(menu_id) AND position >= sqlc.arg(position);

-- name: MenuItemGetNextMainItemPosition :one
SELECT position FROM menu_items WHERE position >= sqlc.arg(position) AND is_main = 1 AND menu_id = sqlc.arg(menu_id) ORDER BY position ASC LIMIT 1;

-- name: CreateMenuItemPost :one
INSERT INTO menu_item_posts (menu_item_id, post_id, position)
VALUES (sqlc.arg(menu_item_id), sqlc.arg(post_id), sqlc.arg(position))
RETURNING *;

-- name: UpdateMenuItemPost :one
UPDATE menu_item_posts
SET menu_item_id = COALESCE(sqlc.narg(menu_item_id), menu_item_id),
    post_id = COALESCE(sqlc.narg(post_id), post_id),
    position = COALESCE(sqlc.narg(position), position)
WHERE menu_item_id = sqlc.narg(menu_item_id) AND post_id = sqlc.arg(post_id)
RETURNING *;

-- name: DeleteMenuItemPost :exec
DELETE FROM menu_item_posts WHERE menu_item_id = sqlc.arg(menu_item_id) AND post_id = sqlc.arg(post_id);

-- name: GetMenuItemPost :one
SELECT * FROM menu_item_posts WHERE menu_item_id = sqlc.arg(menu_item_id) AND post_id = sqlc.arg(post_id);
