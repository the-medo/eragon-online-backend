-- name: CreateMap :one
INSERT INTO maps (name, type, description, width, height, thumbnail_image_id)
VALUES (sqlc.arg(name), sqlc.narg(type), sqlc.narg(description), sqlc.arg(width), sqlc.arg(height), sqlc.narg(thumbnail_image_id))
RETURNING *;

-- name: GetWorldMaps :many
SELECT
    vm.*
FROM
    view_maps vm
    JOIN world_maps wm ON wm.map_id = vm.id
WHERE
    wm.world_id = sqlc.arg(world_id);
;

-- name: GetMapByID :one
SELECT * FROM view_maps WHERE id = sqlc.arg(id);

-- name: UpdateMap :one
UPDATE maps
SET
    name = COALESCE(sqlc.narg(name), name),
    type = COALESCE(sqlc.narg(type), type),
    description = COALESCE(sqlc.narg(description), description),
    width = COALESCE(sqlc.narg(width), width),
    height = COALESCE(sqlc.narg(height), height),
    thumbnail_image_id = COALESCE(sqlc.narg(thumbnail_image_id), thumbnail_image_id)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteMap :exec
DELETE FROM maps WHERE id = sqlc.arg(id);

-- name: CreateMapLayer :one
INSERT INTO map_layers (name, map_id, image_id, is_main, enabled, sublayer)
VALUES (sqlc.arg(name), sqlc.arg(map_id), sqlc.arg(image_id), sqlc.arg(is_main), sqlc.arg(enabled), sqlc.arg(sublayer))
RETURNING *;

-- name: GetMapLayers :many
SELECT * FROM view_map_layers WHERE map_id = sqlc.arg(map_id);

-- name: GetMapLayerByID :one
SELECT * FROM view_map_layers WHERE id = sqlc.arg(map_layer_id);

-- name: UpdateMapLayer :one
UPDATE map_layers
SET
    name = COALESCE(sqlc.narg(name), name),
    image_id = COALESCE(sqlc.narg(image_id), image_id),
    enabled = COALESCE(sqlc.narg(enabled), enabled),
    sublayer = COALESCE(sqlc.narg(sublayer), sublayer)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateMapLayerIsMain :exec
CALL update_map_layer_is_main(sqlc.arg(map_layer_id));

-- name: DeleteMapLayer :exec
DELETE FROM map_layers WHERE id = sqlc.arg(id);

-- name: DeleteMapLayersForMap :exec
DELETE FROM map_layers WHERE map_id = sqlc.arg(map_id);

-- name: CreateWorldMap :one
INSERT INTO world_maps (world_id, map_id)
VALUES (sqlc.arg(world_id), sqlc.arg(map_id))
RETURNING *;

-- name: DeleteWorldMap :exec
DELETE FROM world_maps
WHERE world_id = sqlc.arg(world_id) AND map_id = sqlc.arg(map_id);



--------------------------------------

-- name: CreateMapPinType :one
INSERT INTO map_pin_types (map_id, shape, background_color, border_color, icon_color, icon, icon_size, width)
VALUES (sqlc.arg(map_id), sqlc.arg(shape), sqlc.arg(background_color), sqlc.arg(border_color), sqlc.arg(icon_color), sqlc.arg(icon), sqlc.arg(icon_size), sqlc.arg(width))
RETURNING *;

-- name: GetMapPinTypesForMap :many
SELECT * FROM map_pin_types WHERE map_id = sqlc.arg(map_id);

-- name: UpdateMapPinType :one
UPDATE map_pin_types
SET
    shape = COALESCE(sqlc.narg(shape), shape),
    background_color = COALESCE(sqlc.narg(background_color), background_color),
    border_color = COALESCE(sqlc.narg(border_color), border_color),
    icon_color = COALESCE(sqlc.narg(icon_color), icon_color),
    icon = COALESCE(sqlc.narg(icon), icon),
    icon_size = COALESCE(sqlc.narg(icon_size), icon_size),
    width = COALESCE(sqlc.narg(width), width)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteMapPinType :exec
DELETE FROM map_pin_types WHERE id = sqlc.arg(id);

-- name: CreateMapPin :one
INSERT INTO map_pins (name, map_id, map_pin_type_id, location_id, map_layer_id, x, y)
VALUES (sqlc.arg(name), sqlc.arg(map_id), sqlc.narg(map_pin_type_id), sqlc.narg(location_id), sqlc.narg(map_layer_id), sqlc.arg(x), sqlc.arg(y))
RETURNING *;

-- name: GetMapPins :many
SELECT * FROM view_map_pins WHERE map_id = sqlc.arg(map_id);

-- name: GetMapPinByID :one
SELECT * FROM view_map_pins WHERE id = sqlc.arg(id);

-- name: UpdateMapPin :one
UPDATE map_pins
SET
    name = COALESCE(sqlc.narg(name), name),
    map_pin_type_id = COALESCE(sqlc.narg(map_pin_type_id), map_pin_type_id),
    location_id = COALESCE(sqlc.narg(location_id), location_id),
    map_layer_id = COALESCE(sqlc.narg(map_layer_id), map_layer_id),
    x = COALESCE(sqlc.narg(x), x),
    y = COALESCE(sqlc.narg(y), y)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteMapPin :exec
DELETE FROM map_pins WHERE id = sqlc.arg(id);

-- name: DeleteMapPinsForMapLayer :exec
DELETE FROM map_pins WHERE map_layer_id = sqlc.arg(map_layer_id);

-- name: DeleteMapPinsForMap :exec
DELETE FROM map_pins WHERE map_id = sqlc.arg(map_id);

-- name: GetMapAssignments :one
SELECT
    CAST(MAX(COALESCE(wl.world_id, 0)) as integer) AS world_id,
    0 AS quest_id
FROM
    maps m
    LEFT JOIN world_maps wm ON m.id = wm.location_id
WHERE m.id = sqlc.arg(map_id)
GROUP BY m.id;
