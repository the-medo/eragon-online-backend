CREATE VIEW view_modules AS
SELECT m.id as module_id,
       m.world_id as module_world_id,
       m.system_id as module_system_id,
       m.character_id as module_character_id,
       m.quest_id as module_quest_id,
       m.module_type as module_type,
       m.menu_id as menu_id,
       m.header_img_id as header_img_id,
       m.thumbnail_img_id as thumbnail_img_id,
       m.avatar_img_id as avatar_img_id,
       i_header.url as image_header,
       i_thumbnail.url as image_thumbnail,
       i_avatar.url as image_avatar,
       tags.tags AS tags
FROM
    modules m
    LEFT JOIN (
        SELECT
            mt.module_id,
            cast(array_agg(mt.tag_id) as integer[]) AS tags
        FROM
            module_tags mt
        GROUP BY mt.module_id
    ) tags ON tags.module_id = m.id
    LEFT JOIN images i_header on m.header_img_id = i_header.id
    LEFT JOIN images i_thumbnail on m.thumbnail_img_id = i_thumbnail.id
    LEFT JOIN images i_avatar on m.avatar_img_id = i_avatar.id
;