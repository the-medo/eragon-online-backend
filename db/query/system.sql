-- name: CreateSystem :one
INSERT INTO systems (
    name,
    based_on,
    short_description
) VALUES (
    @name, @based_on, @short_description
) RETURNING *;

-- name: UpdateSystem :one
UPDATE systems
SET
    name = COALESCE(sqlc.narg(name), name),
    based_on = COALESCE(sqlc.narg(based_on), based_on),
    public = COALESCE(sqlc.narg(public), public),
    short_description = COALESCE(sqlc.narg(short_description), short_description)
WHERE
    id = sqlc.arg(system_id)
RETURNING *;

-- name: DeleteSystem :exec
DELETE FROM systems WHERE id = @system_id;

-- name: GetSystemByID :one
SELECT * FROM systems WHERE id = @system_id LIMIT 1;

-- name: GetSystemsByIDs :many
SELECT * FROM systems WHERE id = ANY(@system_ids::int[]);

-- name: GetSystems :many
SELECT * FROM get_systems(@is_public::boolean, @tags::integer[], @order_by::VARCHAR, 'DESC', @page_limit, @page_offset);

-- name: GetSystemsCount :one
SELECT COUNT(*) FROM view_systems
WHERE (@is_public::boolean IS NULL OR public = @is_public) AND
    (array_length(@tags::integer[], 1) IS NULL OR tags @> @tags::integer[]);
