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