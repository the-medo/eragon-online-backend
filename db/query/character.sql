-- name: CreateCharacter :one
INSERT INTO characters (
    name,
    short_description,
    world_id,
    system_id
) VALUES (
    @name, @short_description, @world_id, @system_id
) RETURNING *;

-- name: UpdateCharacter :one
UPDATE characters
SET
    name = COALESCE(sqlc.narg(name), name),
    short_description = COALESCE(sqlc.narg(short_description), short_description),
    public = COALESCE(sqlc.narg(public), public),
    world_id = COALESCE(sqlc.narg(world_id), world_id),
    system_id = COALESCE(sqlc.narg(system_id), system_id)
WHERE
    id = sqlc.arg(character_id)
    RETURNING *;

-- name: DeleteCharacter :exec
DELETE FROM characters WHERE id = @character_id;

-- name: GetCharacterByID :one
SELECT * FROM characters WHERE id = @character_id LIMIT 1;

-- name: GetCharactersByIDs :many
SELECT * FROM characters WHERE id = ANY(@character_ids::int[]);

-- name: GetCharacters :many
SELECT * FROM get_characters(@is_public::boolean, @tags::integer[], @world_id::integer, @system_id::integer, @order_by::VARCHAR, 'DESC', @page_limit, @page_offset);

-- name: GetCharactersCount :one
SELECT COUNT(*) FROM view_characters
WHERE (@is_public::boolean IS NULL OR public = @is_public) AND
    (array_length(@tags::integer[], 1) IS NULL OR tags @> @tags::integer[]);
