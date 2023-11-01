-- name: CreateWorld :one
INSERT INTO worlds (
    name,
    based_on,
    short_description
) VALUES (
     @name, @based_on, @short_description
 ) RETURNING *;

-- name: UpdateWorld :one
UPDATE worlds
SET
    name = COALESCE(sqlc.narg(name), name),
    based_on = COALESCE(sqlc.narg(based_on), based_on),
    public = COALESCE(sqlc.narg(public), public),
    short_description = COALESCE(sqlc.narg(short_description), short_description),
    description_post_id = COALESCE(sqlc.narg(description_post_id), description_post_id)
WHERE
    id = sqlc.arg(world_id)
RETURNING *;

-- name: DeleteWorld :exec
DELETE FROM worlds WHERE id = @world_id;

-- name: GetWorldByID :one
SELECT * FROM view_worlds WHERE id = @world_id LIMIT 1;

-- name: GetWorlds :many
SELECT * FROM get_worlds(@is_public::boolean, @tags::integer[], @order_by::VARCHAR, 'DESC', @page_limit, @page_offset);

-- name: GetWorldsCount :one
SELECT COUNT(*) FROM view_worlds
WHERE (@is_public::boolean IS NULL OR public = @is_public) AND
    (array_length(@tags::integer[], 1) IS NULL OR tags @> @tags::integer[]);

-- name: InsertWorldAdmin :one
INSERT INTO world_admins (
    world_id,
    user_id,
    super_admin,
    approved,
    motivational_letter
) VALUES (@world_id, @user_id, @super_admin, @approved, @motivational_letter) RETURNING *;


-- name: UpdateWorldAdmin :one
UPDATE world_admins
SET
    super_admin = COALESCE(sqlc.narg(super_admin), super_admin),
    approved = COALESCE(sqlc.narg(approved), approved),
    motivational_letter = COALESCE(sqlc.narg(motivational_letter), motivational_letter)
WHERE
    world_id = sqlc.arg(world_id) AND user_id = sqlc.arg(user_id)
RETURNING *;

-- name: DeleteWorldAdmin :exec
DELETE FROM world_admins WHERE world_id = sqlc.arg(world_id) AND user_id = sqlc.arg(user_id);

-- name: GetWorldsOfUser :many
SELECT
    vw.*,
    1 as world_admin,
    wa.super_admin as world_super_admin
FROM
    view_worlds vw
    JOIN world_admins wa ON wa.world_id = vw.id
WHERE
    wa.user_id = @user_id AND wa.approved = 1
;

-- name: GetWorldAdmins :many
SELECT
    vu.*,
    wa.world_id as world_id,
    wa.created_at as world_admin_created_at,
    wa.super_admin as world_admin_super_admin,
    wa.approved as world_admin_approved,
    wa.motivational_letter as world_admin_motivational_letter
FROM
    view_users vu
    JOIN world_admins wa on wa.user_id = vu.id
WHERE
    wa.world_id = sqlc.arg(world_id)
;

-- name: IsWorldAdmin :one
SELECT * FROM world_admins WHERE user_id = @user_id AND world_id = @world_id AND approved = 1;

-- name: IsWorldSuperAdmin :one
SELECT * FROM world_admins WHERE user_id = @user_id AND world_id = @world_id AND approved = 1 AND super_admin = 1;

-- name: CreateModuleMapPinTypeGroup :one
INSERT INTO module_map_pin_type_groups (module_id, map_pin_type_group_id)
VALUES (sqlc.arg(module_id), sqlc.arg(map_pin_type_group_id))
ON CONFLICT (module_id, map_pin_type_group_id) DO NOTHING
RETURNING *;

-- name: DeleteModuleMapPinTypeGroup :exec
DELETE FROM module_map_pin_type_groups WHERE module_id = sqlc.arg(module_id) AND map_pin_type_group_id = sqlc.arg(map_pin_type_group_id);