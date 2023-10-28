-- name: GetModuleId :one
SELECT
    id as module_id, module_type
FROM modules
WHERE
    world_id = COALESCE(sqlc.narg(world_id), world_id) AND
    quest_id = COALESCE(sqlc.narg(quest_id), quest_id) AND
    character_id = COALESCE(sqlc.narg(character_id), character_id) AND
    system_id = COALESCE(sqlc.narg(system_id), system_id)
;