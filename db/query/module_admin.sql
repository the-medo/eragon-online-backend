
-- name: CreateModuleAdmin :one
INSERT INTO module_admins (
    module_id,
    user_id,
    super_admin,
    approved,
    motivational_letter,
    allowed_entity_types,
    allowed_menu
) VALUES (@module_id, @user_id, @super_admin, @approved, @motivational_letter, @allowed_entity_types, @allowed_menu) RETURNING *;


-- name: UpdateModuleAdmin :one
UPDATE module_admins
SET
    super_admin = COALESCE(sqlc.narg(super_admin), super_admin),
    approved = COALESCE(sqlc.narg(approved), approved),
    motivational_letter = COALESCE(sqlc.narg(motivational_letter), motivational_letter),
    allowed_entity_types = COALESCE(sqlc.narg(allowed_entity_types), allowed_entity_types),
    allowed_menu = COALESCE(sqlc.narg(allowed_menu), allowed_menu)
WHERE
    module_id = sqlc.arg(module_id) AND user_id = sqlc.arg(user_id)
RETURNING *;

-- name: DeleteModuleAdmin :exec
DELETE FROM module_admins WHERE module_id = sqlc.arg(module_id) AND user_id = sqlc.arg(user_id);

-- name: GetModulesOfAdmin :many
SELECT
    *
FROM
    view_module_admins
WHERE
    user_id = @user_id AND approved = 1
;

-- name: GetModuleAdmins :many
SELECT
    vu.*,
    ma.module_id           as module_id,
    ma.created_at          as module_admin_created_at,
    ma.super_admin         as module_admin_super_admin,
    ma.approved            as module_admin_approved,
    ma.motivational_letter as module_admin_motivational_letter,
    ma.allowed_entity_types as module_admin_allowed_entity_types,
    ma.allowed_menu        as module_admin_allowed_menu
FROM
    view_users vu
    JOIN module_admins ma on ma.user_id = vu.id
WHERE
    ma.module_id = sqlc.arg(module_id)
;

-- name: GetModuleAdmin :one
SELECT * FROM view_module_admins WHERE user_id = @user_id AND id = @module_id;

-- name: GetModuleAdminByMenuId :one
SELECT * FROM view_module_admins WHERE user_id = sqlc.arg(user_id) AND menu_id = sqlc.arg(menu_id);

-- name: GetEntityModuleAdmin :one
SELECT
    vma.*
FROM
    entities e
    JOIN view_module_admins vma ON e.module_id = vma.id
WHERE
    e.id = sqlc.arg(entity_id) AND
    vma.user_id = sqlc.arg(user_id) AND
    (e.type = ANY(vma.allowed_entity_types) OR vma.super_admin = true)
;