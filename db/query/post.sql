-- name: CreatePost :one
INSERT INTO posts
(
    user_id,
    title,
    content,
    description,
    is_draft,
    is_private,
    thumbnail_img_id
)
VALUES
    (sqlc.arg(user_id), sqlc.arg(title), sqlc.arg(content), sqlc.arg(description), sqlc.arg(is_draft), sqlc.arg(is_private), sqlc.arg(thumbnail_img_id))
RETURNING *;

-- name: GetPostById :one
SELECT
    *
FROM
    posts
WHERE
    id = sqlc.arg(post_id);

-- name: GetPostsByIDs :many
SELECT * FROM posts WHERE id = ANY(@post_ids::int[]);

-- name: GetPostsByUserId :many
SELECT
    *
FROM
    view_posts
WHERE
    user_id = sqlc.arg(user_id) AND
    deleted_at IS NULL
ORDER BY created_at DESC
LIMIT @page_limit
OFFSET @page_offset;

-- name: GetPostsByModule :many
WITH cte AS (
    SELECT
        p.*
    FROM
        posts p
        LEFT JOIN view_entities e ON e.post_id = p.id
        LEFT JOIN modules m ON m.id = e.module_id
    WHERE
        m.id = sqlc.arg(module_id) AND
        p.deleted_at IS NULL
)
SELECT
    CAST((SELECT count(*) FROM cte) as integer) as total_count,
    cte.*
FROM cte
ORDER BY created_at DESC
LIMIT @page_limit
OFFSET @page_offset;

-- name: UpdatePost :one
UPDATE posts
SET
    title = COALESCE(sqlc.narg(title), title),
    content = COALESCE(sqlc.narg(content), content),
    description = COALESCE(sqlc.narg(description), description),
    is_draft = COALESCE(sqlc.narg(is_draft), is_draft),
    is_private = COALESCE(sqlc.narg(is_private), is_private),
    last_updated_user_id = sqlc.arg(last_updated_user_id),
    last_updated_at = now(),
    thumbnail_img_id = COALESCE(sqlc.narg(thumbnail_img_id), thumbnail_img_id)
WHERE
    posts.id = sqlc.arg(post_id)
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
    user_id,
    title,
    content,
    created_at,
    deleted_at,
    last_updated_at,
    last_updated_user_id,
    is_draft,
    is_private,
    description,
    thumbnail_img_id
)
SELECT
    id,
    user_id,
    title,
    content,
    created_at,
    deleted_at,
    last_updated_at,
    last_updated_user_id,
    is_draft,
    is_private,
    description,
    thumbnail_img_id
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
    last_updated_user_id,
    is_draft,
    is_private,
    description,
    thumbnail_img_id
FROM post_history WHERE post_id = sqlc.arg(post_id) ORDER BY created_at DESC;

-- name: GetPostHistoryById :one
SELECT
    id as post_history_id,
    post_id,
    user_id,
    title,
    content,
    created_at,
    deleted_at,
    last_updated_at,
    last_updated_user_id,
    is_draft,
    is_private,
    description,
    thumbnail_img_id
FROM post_history WHERE id = sqlc.arg(post_history_id);
