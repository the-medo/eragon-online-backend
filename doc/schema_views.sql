CREATE VIEW view_worlds AS
SELECT
    w.*,
    i_avatar.url as image_avatar,
    i_header.url as image_header,
    ws.final_content_rating as rating,
    ws.final_activity as activity
FROM
    worlds w
    JOIN world_images wi ON w.id = wi.world_id
    JOIN world_stats ws ON w.id = wi.world_id
    LEFT JOIN images i_avatar on wi.image_avatar = i_avatar.id
    LEFT JOIN images i_header on wi.image_header = i_header.id
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