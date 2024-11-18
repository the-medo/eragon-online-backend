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
SELECT * FROM get_quests(@is_public::boolean, @tags::integer[], @world_id::integer, @system_id::integer, @order_by::VARCHAR, 'DESC', @page_limit, @page_offset);

-- name: GetQuestsCount :one
SELECT COUNT(*) FROM view_quests
WHERE (@is_public::boolean IS NULL OR public = @is_public) AND
    (array_length(@tags::integer[], 1) IS NULL OR tags @> @tags::integer[]);

-- name: CreateQuestSetting :one
INSERT INTO quest_settings (quest_id, status, can_join) VALUES (@quest_id, @status, @can_join) RETURNING *;

-- name: CreateQuestCharacter :one
INSERT INTO quest_characters (quest_id, character_id, created_at, approved, motivational_letter)
VALUES (@quest_id, @character_id, NOW(), @approved, @motivational_letter) RETURNING *;

-- name: UpdateQuestSetting :one
UPDATE quest_settings
SET
    status = COALESCE(@status, status),
    can_join = COALESCE(@can_join, can_join)
WHERE
    quest_id = @quest_id
RETURNING *;

-- name: DeleteQuestSetting :exec
DELETE FROM quest_settings
WHERE quest_id = @quest_id;

-- name: GetQuestSettingByQuestID :one
SELECT * FROM quest_settings
WHERE quest_id = @quest_id
LIMIT 1;

-- name: UpdateQuestCharacter :one
UPDATE quest_characters
SET
    approved = COALESCE(@approved, approved),
    motivational_letter = COALESCE(@motivational_letter, motivational_letter)
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

