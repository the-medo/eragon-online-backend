-- name: CreatePost :one
INSERT INTO posts
(
    user_id,
    title,
    content
)
VALUES
    (sqlc.arg(user_id), sqlc.arg(title), sqlc.arg(content))
RETURNING *;

-- name: GetPostById :one
SELECT * FROM posts WHERE id = sqlc.arg(post_id);

-- name: GetPostsByUserId :many
SELECT * FROM posts WHERE user_id = sqlc.arg(user_id) AND deleted_at IS NULL ORDER BY created_at DESC;

-- name: UpdatePost :one
UPDATE posts
SET
    title = COALESCE(sqlc.arg(title), title),
    content = COALESCE(sqlc.arg(content), content),
    last_updated_at = now(),
    last_updated_user_id = sqlc.arg(last_updated_user_id)
WHERE
    id = sqlc.arg(post_id)
RETURNING *;

-- name: DeletePost :exec
UPDATE posts
SET
    deleted_at = now()
WHERE
    id = sqlc.arg(id);

-- name: InsertPostHistory :one
INSERT INTO post_history (
    post_id,
    user_id,
    title,
    content,
    created_at,
    deleted_at,
    last_updated_at,
    last_updated_user_id
)
SELECT
    id,
    user_id,
    title,
    content,
    created_at,
    deleted_at,
    last_updated_at,
    last_updated_user_id
FROM
    posts
WHERE
    posts.id = sqlc.arg(post_id)
RETURNING *;

-- name: GetPostHistoryByPostId :many
SELECT
    id as post_history_id,
    post_id,
    user_id,
    title,
    created_at,
    deleted_at,
    last_updated_at,
    last_updated_user_id
FROM post_history WHERE post_id = sqlc.arg(post_id) ORDER BY created_at DESC;

-- name: GetPostHistoryById :many
SELECT
    id as post_history_id,
    post_id,
    user_id,
    title,
    content,
    created_at,
    deleted_at,
    last_updated_at,
    last_updated_user_id
FROM post_history WHERE id = sqlc.arg(post_history_id);