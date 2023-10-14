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
WITH deleted_item AS (
    DELETE FROM "menu_items" d
        WHERE d.id = sqlc.arg(menu_item_id)
        RETURNING *
)
UPDATE "menu_items"
SET "position" = "position" - 1
WHERE "menu_id" = (SELECT menu_id FROM deleted_item)
  AND "position" > (SELECT position FROM deleted_item);

-- name: GetMenuItems :many
SELECT * FROM menu_items WHERE menu_id = sqlc.arg(menu_id);

-- name: GetMenuItemById :one
SELECT * FROM menu_items WHERE id = sqlc.arg(id);

-- name: MenuItemPostChangePositions :exec
CALL move_menu_item_post(sqlc.arg(menu_item_id), sqlc.arg(post_id), sqlc.arg(target_position));

-- name: MenuItemChangePositions :exec
CALL move_menu_item(sqlc.arg(menu_item_id), sqlc.arg(target_position));

-- name: MenuItemMoveGroupUp :exec
CALL move_group_up(sqlc.arg(menu_item_id));

-- name: MenuEntityGroupChangePositions :exec
CALL move_menu_entity_groups(sqlc.arg(menu_id), sqlc.arg(entity_group_id), sqlc.arg(target_position));

-- name: CreateMenuItemPost :one
WITH post_count AS (
    SELECT COUNT(*) AS count FROM menu_item_posts WHERE menu_item_id = COALESCE(sqlc.narg(menu_item_id), 0)
)
INSERT INTO menu_item_posts (menu_id, menu_item_id, post_id, position)
SELECT
    sqlc.arg(menu_id) as menu_id,
    sqlc.narg(menu_item_id) as menu_item_id,
    sqlc.arg(post_id) as post_id,
    COALESCE(sqlc.narg(position), count + 1) as position
FROM post_count
RETURNING *;

-- name: UpdateMenuItemPost :one
UPDATE menu_item_posts
SET menu_item_id = COALESCE(sqlc.narg(new_menu_item_id), menu_item_id),
    post_id = COALESCE(sqlc.narg(post_id), post_id),
    position = COALESCE(sqlc.narg(position), position)
WHERE (menu_item_id = sqlc.narg(menu_item_id) OR (sqlc.narg(menu_item_id) IS NULL AND menu_item_id IS NULL))  AND post_id = sqlc.arg(post_id)
RETURNING *;

-- name: UnassignMenuItemPost :one
UPDATE menu_item_posts
SET menu_item_id = NULL
WHERE menu_item_id = sqlc.arg(menu_item_id) AND post_id = sqlc.arg(post_id)
RETURNING *;

-- name: DeleteMenuItemPost :exec
WITH deleted_menu_item_post AS (
    DELETE FROM "menu_item_posts" d
        WHERE d.menu_id = sqlc.arg(menu_id) AND d.post_id = sqlc.arg(post_id)
        RETURNING *
)
UPDATE "menu_item_posts"
SET "position" = "position" - 1
WHERE
    "menu_item_id" = (SELECT menu_item_id FROM deleted_menu_item_post)
    AND "position" > (SELECT position FROM deleted_menu_item_post);

-- name: GetMenuItemPost :one
SELECT * FROM view_menu_item_posts WHERE menu_item_id = sqlc.arg(menu_item_id) AND post_id = sqlc.arg(post_id);

-- name: GetMenuItemPosts :many
SELECT * FROM view_menu_item_posts WHERE menu_item_id = sqlc.arg(menu_item_id);

-- name: GetMenuItemPostsByMenuId :many
SELECT * FROM view_menu_item_posts WHERE menu_id = sqlc.arg(menu_id);