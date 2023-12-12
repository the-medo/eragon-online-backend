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
            wt.world_id,
            cast(array_agg(t.tag) as varchar[]) AS tags
        FROM
            world_tags wt
                LEFT JOIN world_tags_available t ON t.id = wt.tag_id
        GROUP BY wt.world_id
    ) tags ON tags.world_id = w.id
    LEFT JOIN images i_header on wi.header_img_id = i_header.id
    LEFT JOIN images i_thumbnail on wi.thumbnail_img_id = i_thumbnail.id
    LEFT JOIN images i_avatar on wi.avatar_img_id = i_avatar.id
;