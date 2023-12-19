CREATE VIEW view_modules AS
SELECT m.id as id,
       m.world_id as world_id,
       m.system_id as system_id,
       m.character_id as character_id,
       m.quest_id as quest_id,
       m.module_type as module_type,
       m.menu_id as menu_id,
       m.header_img_id as header_img_id,
       m.thumbnail_img_id as thumbnail_img_id,
       m.avatar_img_id as avatar_img_id,
       m.description_post_id as description_post_id,
       cast(array_agg(tags.tag_id) as integer[]) AS tags
FROM
    modules m
        LEFT JOIN module_tags tags ON tags.module_id = m.id
GROUP BY m.id
;