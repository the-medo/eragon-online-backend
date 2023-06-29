-- name: CreatePost :one
INSERT INTO posts
(
    user_id,
    title,
    post_type_id,
    content,
    is_draft,
    is_private
)
VALUES
    (sqlc.arg(user_id), sqlc.arg(title), sqlc.arg(post_type_id), sqlc.arg(content), sqlc.arg(is_draft), sqlc.arg(is_private))
RETURNING *;

-- name: GetPostById :one
SELECT
    *
FROM
    view_posts
WHERE
    id = sqlc.arg(post_id);

-- name: GetPostTypeById :one
SELECT * FROM post_types WHERE id = sqlc.arg(post_type_id);

-- name: GetPostsByUserId :many
SELECT
    *
FROM
    view_posts
WHERE
    user_id = sqlc.arg(user_id) AND
    post_type_id = COALESCE(sqlc.narg(post_type_id), post_type_id) AND
    deleted_at IS NULL
ORDER BY created_at DESC
LIMIT @page_limit
OFFSET @page_offset;

-- name: UpdatePost :one
UPDATE posts
SET
    title = COALESCE(sqlc.narg(title), title),
    content = COALESCE(sqlc.narg(content), content),
    post_type_id = COALESCE(sqlc.narg(post_type_id), post_type_id),
    is_draft = COALESCE(sqlc.narg(is_draft), is_draft),
    is_private = COALESCE(sqlc.narg(is_private), is_private),
    last_updated_user_id = sqlc.arg(last_updated_user_id),
    last_updated_at = now()
WHERE
    id = sqlc.arg(post_id)
RETURNING *;

-- name: DeletePost :exec
UPDATE posts
SET
    deleted_at = now()
WHERE
    id = sqlc.arg(post_id);

-- name: InsertPostHistory :one
INSERT INTO post_history (
    post_id,
    post_type_id,
    user_id,
    title,
    content,
    created_at,
    deleted_at,
    last_updated_at,
    last_updated_user_id,
    is_draft,
    is_private
)
SELECT
    id,
    post_type_id,
    user_id,
    title,
    content,
    created_at,
    deleted_at,
    last_updated_at,
    last_updated_user_id,
    is_draft,
    is_private
FROM
    posts
WHERE
    posts.id = sqlc.arg(post_id)
RETURNING *;

-- name: GetPostHistoryByPostId :many
SELECT
    id as post_history_id,
    post_id,
    post_type_id,
    user_id,
    title,
    created_at,
    deleted_at,
    last_updated_at,
    last_updated_user_id,
    is_draft,
    is_private
FROM post_history WHERE post_id = sqlc.arg(post_id) ORDER BY created_at DESC;

-- name: GetPostHistoryById :one
SELECT
    id as post_history_id,
    post_id,
    post_type_id,
    user_id,
    title,
    content,
    created_at,
    deleted_at,
    last_updated_at,
    last_updated_user_id,
    is_draft,
    is_private
FROM post_history WHERE id = sqlc.arg(post_history_id);