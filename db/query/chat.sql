-- name: AddChatMessage :one
INSERT INTO chat
(
    user_id,
    text
)
VALUES (@user_id, @text)
RETURNING *;

-- name: GetChatMessages :many
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
LIMIT @page_limit
OFFSET @page_offset;

-- name: GetChatMessage :one
SELECT
    c.id as chat_id,
    c.text as text,
    c.created_at as created_at,
    c.user_id as user_id,
    u.username as username
FROM
    chat c
    JOIN users u ON c.user_id = u.id
WHERE c.id = @id;

-- name: DeleteChatMessage :exec
DELETE FROM chat WHERE id = @id;