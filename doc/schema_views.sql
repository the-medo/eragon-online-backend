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