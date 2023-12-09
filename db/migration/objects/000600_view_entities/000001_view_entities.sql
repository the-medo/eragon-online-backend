CREATE VIEW view_entities AS
SELECT
    e.*,
    m.module_type as module_type,
    CAST(CASE m.module_type
         WHEN 'world' THEN m.world_id
         WHEN 'quest' THEN m.quest_id
         WHEN 'character' THEN m.character_id
         WHEN 'system' THEN m.system_id
    END as integer) as module_type_id,
    tags.tags as tags
FROM
    entities e
    LEFT JOIN modules m ON e.module_id = m.id
    LEFT JOIN (
        SELECT
            et.entity_id,
            cast(array_agg(tag_available.id) as int[]) AS tags
        FROM
            entity_tags et
                LEFT JOIN module_entity_tags_available tag_available ON tag_available.id = et.tag_id
        GROUP BY et.entity_id
    ) tags ON tags.entity_id = e.id
;