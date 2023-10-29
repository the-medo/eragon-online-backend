-- name: GetModuleId :one
SELECT
    id as module_id, module_type
FROM modules
WHERE
    world_id = COALESCE(sqlc.narg(world_id), world_id) OR
    quest_id = COALESCE(sqlc.narg(quest_id), quest_id) OR
    character_id = COALESCE(sqlc.narg(character_id), character_id) OR
    system_id = COALESCE(sqlc.narg(system_id), system_id)
;