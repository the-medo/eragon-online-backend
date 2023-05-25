-- name: CreateUser :one
INSERT INTO users
(
    username,
    hashed_password,
    email
)
VALUES
    ($1, $2, $3)
RETURNING *;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
    hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
    password_changed_at = COALESCE(sqlc.narg(password_changed_at), password_changed_at),
    username = COALESCE(sqlc.narg(username), username),
    email = COALESCE(sqlc.narg(email), email),
    img_id = COALESCE(sqlc.narg(img_id), img_id),
    is_email_verified = COALESCE(sqlc.narg(is_email_verified), is_email_verified)
WHERE
    id = sqlc.arg(id)
RETURNING *;

-- name: HasUserRole :one
SELECT
    *
FROM
    user_roles ur
    JOIN roles r on ur.role_id = r.id
WHERE user_id = @user_id AND r.name = @role LIMIT 1;

-- name: GetUserRoles :many
SELECT
    ur.*,
    r.name AS role_name,
    r.description AS role_description
FROM
    user_roles ur
    JOIN roles r ON ur.role_id = r.id
WHERE user_id = @user_id;

-- name: AddUserRole :one
INSERT INTO user_roles (user_id, role_id) VALUES (@user_id, @role_id) RETURNING *;

-- name: RemoveUserRole :exec
DELETE FROM user_roles WHERE user_id = @user_id AND role_id = @role_id;

-- name: GetUsers :many
SELECT
    *
FROM
    users AS u
    LEFT JOIN images i ON u.img_id = i.id
ORDER BY username
LIMIT @page_limit
OFFSET @page_offset
;


-- name: AddUserPasswordReset :one
INSERT INTO user_password_reset (user_id, code) VALUES (@user_id, @code) RETURNING *;

-- name: GetUserPasswordReset :one
SELECT * FROM user_password_reset WHERE code = @code AND expired_at > NOW() LIMIT 1;

-- name: DeleteUserPasswordReset :exec
DELETE FROM user_password_reset WHERE user_id = @user_id AND code = @code;