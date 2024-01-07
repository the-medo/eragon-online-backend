-- name: CreateMenu :one
INSERT INTO menus (menu_code, menu_header_img_id)
VALUES (sqlc.arg(menu_code), sqlc.narg(menu_header_img_id))
RETURNING *;

-- name: UpdateMenu :one
UPDATE menus
SET menu_code = COALESCE(sqlc.narg(menu_code), menu_code),
    menu_header_img_id = COALESCE(sqlc.narg(menu_header_img_id), menu_header_img_id)
WHERE menus.id = sqlc.arg(id)
RETURNING *
;

-- name: DeleteMenu :exec
DELETE FROM menus WHERE id = sqlc.arg(id);

-- name: GetMenu :one
SELECT * FROM view_menus WHERE id = sqlc.arg(id);

-- name: CreateMenuItem :one
INSERT INTO menu_items (menu_id, menu_item_code, name, position, is_main, description_post_id, entity_group_id)
VALUES (sqlc.arg(menu_id), sqlc.arg(menu_item_code), sqlc.arg(name), sqlc.arg(position), sqlc.narg(is_main), sqlc.narg(description_post_id), sqlc.narg(entity_group_id))
RETURNING *;

-- name: UpdateMenuItem :one
UPDATE menu_items
SET
    menu_item_code = COALESCE(sqlc.narg(menu_item_code), menu_item_code),
    name = COALESCE(sqlc.narg(name), name),
    -- position = COALESCE(sqlc.narg(position), position),
    is_main = COALESCE(sqlc.narg(is_main), is_main),
    description_post_id = COALESCE(sqlc.narg(description_post_id), description_post_id),
    entity_group_id = COALESCE(sqlc.narg(entity_group_id), entity_group_id)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteMenuItem :exec
CALL delete_menu_item(sqlc.arg(menu_item_id));

-- name: GetMenuItems :many
SELECT * FROM menu_items WHERE menu_id = sqlc.arg(menu_id);

-- name: GetMenuItemById :one
SELECT * FROM menu_items WHERE id = sqlc.arg(id);

-- name: MenuItemChangePositions :exec
CALL move_menu_item(sqlc.arg(menu_item_id), sqlc.arg(target_position));

-- name: MenuItemMoveGroupUp :exec
CALL move_group_up(sqlc.arg(menu_item_id));

-- name: GetRecursiveEntities :many
SELECT * FROM get_recursive_entities(sqlc.arg(entity_group_id));