-- name: CreateQuest :one
INSERT INTO quests (
    name,
    short_description,
    world_id,
    system_id
) VALUES (
    @name, @short_description, @world_id, @system_id
) RETURNING *;

-- name: UpdateQuest :one
UPDATE quests
SET
    name = COALESCE(sqlc.narg(name), name),
    short_description = COALESCE(sqlc.narg(short_description), short_description),
    public = COALESCE(sqlc.narg(public), public),
    world_id = COALESCE(sqlc.narg(world_id), world_id),
    system_id = COALESCE(sqlc.narg(system_id), system_id)
WHERE
    id = sqlc.arg(quest_id)
    RETURNING *;

-- name: DeleteQuest :exec
DELETE FROM quests WHERE id = @quest_id;

-- name: GetQuestByID :one
SELECT * FROM quests WHERE id = @quest_id LIMIT 1;

-- name: GetQuestsByIDs :many
SELECT * FROM quests WHERE id = ANY(@quest_ids::int[]);

-- name: GetQuests :many
SELECT * FROM get_quests(@is_public::boolean, @tags::integer[], @order_by::VARCHAR, 'DESC', @page_limit, @page_offset);

-- name: GetQuestsCount :one
SELECT COUNT(*) FROM view_quests
WHERE (@is_public::boolean IS NULL OR public = @is_public) AND
    (array_length(@tags::integer[], 1) IS NULL OR tags @> @tags::integer[]);
