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
    system_id = COALESCE(sqlc.narg(system_id), system_id),
    status = COALESCE(sqlc.narg(status), status),
    can_join = COALESCE(sqlc.narg(can_join), can_join)
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
SELECT * FROM get_quests(sqlc.narg(is_public)::boolean, sqlc.narg(tags)::integer[], sqlc.narg(world_id)::integer, sqlc.narg(system_id)::integer, sqlc.narg(can_join)::boolean, sqlc.narg(status)::quest_status, sqlc.narg(order_by)::VARCHAR, 'DESC', sqlc.narg(page_limit), sqlc.narg(page_offset));

-- name: GetQuestsCount :one
SELECT COUNT(*) FROM view_quests
WHERE (sqlc.narg(is_public)::boolean IS NULL OR public = sqlc.narg(is_public)) AND
      (sqlc.narg(world_id)::int IS NULL OR world_id = sqlc.narg(world_id)) AND
      (sqlc.narg(system_id)::int IS NULL OR system_id = sqlc.narg(system_id)) AND
      (sqlc.narg(can_join)::boolean IS NULL OR can_join = sqlc.narg(can_join)) AND
      (sqlc.narg(status)::quest_status IS NULL OR status = sqlc.narg(status)) AND
    (array_length(sqlc.narg(tags)::integer[], 1) IS NULL OR tags @> sqlc.narg(tags)::integer[]);

-- name: CreateQuestCharacter :one
INSERT INTO quest_characters (quest_id, character_id, created_at, approved, motivational_letter)
VALUES (@quest_id, @character_id, NOW(), @approved, @motivational_letter) RETURNING *;

-- name: UpdateQuestCharacter :one
UPDATE quest_characters
SET
    approved = COALESCE(sqlc.narg(approved), approved),
    motivational_letter = COALESCE(sqlc.narg(motivational_letter), motivational_letter)
WHERE
    quest_id = @quest_id AND character_id = @character_id
RETURNING *;

-- name: DeleteQuestCharacter :exec
DELETE FROM quest_characters
WHERE quest_id = @quest_id AND character_id = @character_id;

-- name: GetQuestCharacterByQuestAndCharacterID :one
SELECT * FROM quest_characters
WHERE quest_id = @quest_id AND character_id = @character_id
LIMIT 1;

-- name: GetQuestCharactersByQuestID :many
SELECT * FROM quest_characters
WHERE quest_id = @quest_id;

-- name: GetQuestCharactersByCharacterID :many
SELECT * FROM quest_characters
WHERE character_id = @character_id;

