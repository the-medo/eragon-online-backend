-- name: AddChatPost :one
INSERT INTO chat
(
    user_id,
    text
)
VALUES (sqlc.arg(user_id), sqlc.arg(text))
RETURNING *;

-- name: GetChatPosts :many
SELECT
    c.id as chat_id,
    c.text as text,
    c.created_at as created_at,
    c.user_id as user_id,
    u.username as username
FROM
    chat c
    JOIN users u ON c.user_id = u.id
ORDER BY c.id DESC
LIMIT sqlc.arg(page_limit)
OFFSET sqlc.arg(page_offset);

-- name: GetChatPost :one
SELECT
    c.id as chat_id,
    c.text as text,
    c.created_at as created_at,
    c.user_id as user_id,
    u.username as username
FROM
    chat c
    JOIN users u ON c.user_id = u.id
WHERE c.id = sqlc.arg(id);

-- name: DeleteChatPost :exec
DELETE FROM chat WHERE id = sqlc.arg(id);