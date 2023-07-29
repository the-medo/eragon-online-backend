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
SELECT * FROM view_worlds
WHERE (@is_public::boolean IS NULL OR public = @is_public)
ORDER BY
    CASE
     WHEN @order_result::bool
         THEN @order_by::VARCHAR
     ELSE 'created_at'
     END
DESC
LIMIT @page_limit
OFFSET @page_offset;

-- name: InsertWorldAdmin :one
INSERT INTO world_admins (
    world_id,
    user_id,
    super_admin,
    approved,
    motivational_letter
) VALUES (@world_id, @user_id, @super_admin, @approved, @motivational_letter) RETURNING *;


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

-- name: GetAdminsOfWorld :many
SELECT
    vu.*,
    wa.super_admin as super_admin
FROM
    view_users vu
    JOIN world_admins wa on wa.user_id = vu.id
WHERE
    wa.world_id = @world_id AND
    wa.approved = 1
;

-- name: IsWorldAdmin :one
SELECT * FROM world_admins WHERE user_id = @user_id AND world_id = @world_id AND approved = 1;

-- name: IsWorldSuperAdmin :one
SELECT * FROM world_admins WHERE user_id = @user_id AND world_id = @world_id AND approved = 1 AND super_admin = 1;