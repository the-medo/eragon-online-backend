CREATE VIEW view_worlds AS
SELECT
    w.*,
    i_header.url as image_header,
    i_thumbnail.url as image_thumbnail,
    i_avatar.url as image_avatar,
    tags.tags AS tags,
    COALESCE(activity.activity_post_count, 0) AS activity_post_count,
    COALESCE(activity.activity_quest_count, 0) AS activity_quest_count,
    COALESCE(activity.activity_resource_count, 0) AS activity_resource_count,
    wm.menu_id as world_menu_id
FROM
    worlds w
        JOIN world_images wi ON w.id = wi.world_id
        JOIN world_menu wm ON w.id = wm.world_id
        LEFT JOIN (
        SELECT
            wa.world_id,
            cast(sum(wa.post_count) as integer) AS activity_post_count,
            cast(sum(wa.quest_count) as integer) AS activity_quest_count,
            cast(sum(wa.resource_count) as integer) AS activity_resource_count
        FROM
            world_activity wa
        WHERE
                wa.date >= (now() - interval '30 days')
        GROUP BY wa.world_id
    ) activity ON activity.world_id = w.id
        LEFT JOIN (
        SELECT
            m.world_id,
            cast(array_agg(t.id) as integer[]) AS tags
        FROM
            modules m
            JOIN module_tags mt ON m.id = mt.module_id
            LEFT JOIN module_type_tags_available t ON t.id = mt.tag_id
        WHERE world_id IS NOT NULL
        GROUP BY m.world_id
    ) tags ON tags.world_id = w.id
        LEFT JOIN images i_header on wi.header_img_id = i_header.id
        LEFT JOIN images i_thumbnail on wi.thumbnail_img_id = i_thumbnail.id
        LEFT JOIN images i_avatar on wi.avatar_img_id = i_avatar.id
;

CREATE VIEW view_users AS
SELECT
    u.*,
    i.id as avatar_image_id,
    i.url as avatar_image_url,
    i.img_guid as avatar_image_guid,
    p.deleted_at as introduction_post_deleted_at
FROM
    users AS u
    LEFT JOIN images i ON u.img_id = i.id
    LEFT JOIN posts p ON u.introduction_post_id = p.id
;


CREATE VIEW view_posts AS
SELECT
    p.*,
    pt.name as post_type_name,
    pt.draftable as post_type_draftable,
    pt.privatable as post_type_privatable,
    i.url as thumbnail_img_url,
    e.id as entity_id,
    e.module_id as module_id,
    e.module_type as module_type,
    e.module_type_id as module_type_id,
    e.tags as tags
FROM
    posts p
    JOIN post_types pt ON p.post_type_id = pt.id
    LEFT JOIN images i ON p.thumbnail_img_id = i.id
    LEFT JOIN view_entities e ON e.post_id = p.id
;


CREATE VIEW view_menu_item_posts AS
SELECT * FROM menu_item_posts mip JOIN view_posts vp ON mip.post_id = vp.id;


CREATE VIEW view_locations AS
SELECT
    l.*,
    i.url as thumbnail_image_url,
    p.title as post_title,
    e.id as entity_id,
    e.module_id as module_id,
    e.module_type as module_type,
    e.module_type_id as module_type_id,
    e.tags as tags
FROM
    locations l
    JOIN view_entities e ON e.location_id = l.id
    LEFT JOIN images i ON l.thumbnail_image_id = i.id
    LEFT JOIN posts p ON l.post_id = p.id
;

CREATE VIEW view_maps AS
SELECT
    m.*,
    i.url as thumbnail_image_url
FROM
    maps m
        LEFT JOIN images i ON m.thumbnail_image_id = i.id
;

CREATE VIEW view_map_layers AS
SELECT
    ml.*,
    i.url as image_url
FROM
    map_layers ml
        LEFT JOIN images i ON ml.image_id = i.id
;


CREATE VIEW view_connections_world_posts AS
SELECT
    wl.world_id as world_id,
    l.post_id as post_id,
    'locations' as helper_name,
    l.id as helper_id
FROM
    world_locations wl
    JOIN locations l ON wl.location_id = l.id JOIN posts p ON l.post_id = p.id
WHERE l.post_id IS NOT NULL

UNION ALL

SELECT
    wm.world_id as world_id,
    mip.post_id,
    'menu_item_posts',
    mip.menu_id
FROM
    world_menu wm
    JOIN menu_item_posts mip ON mip.menu_id = wm.menu_id

UNION ALL

SELECT
    wm.world_id,
    mi.description_post_id,
    'menu_items',
    mi.id
FROM
    world_menu wm
    JOIN menu_items mi ON mi.menu_id = wm.menu_id
WHERE mi.description_post_id IS NOT NULL

UNION ALL

SELECT
    id,
    description_post_id,
    'worlds',
    id
FROM
    worlds
WHERE description_post_id IS NOT NULL

UNION ALL

SELECT
    wm.world_id,
    e.post_id,
    'entities',
    e.id
FROM
    world_menu wm
    JOIN menu_items mi ON mi.menu_id = wm.menu_id
    JOIN get_recursive_entities(mi.entity_group_id) re ON 1 = 1
    JOIN entities e ON e.id = re.content_entity_id
WHERE e.post_id IS NOT NULL
;


CREATE VIEW view_connections_menus AS
SELECT
    m.id as menu_id,
    COALESCE(wm.world_id, 0) as world_id,
    0 as quest_id,
    0 as character_id,
    0 as system_id
FROM
    menus m
        LEFT JOIN world_menu wm ON m.id = wm.menu_id
--     LEFT JOIN quest_menu q ON m.id = q.menu_id
--     LEFT JOIN character_menu c ON m.id = c.menu_id
--     LEFT JOIN system_menu s ON m.id = s.menu_id
;