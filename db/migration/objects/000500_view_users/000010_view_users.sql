CREATE VIEW view_users AS
SELECT
    u.*,
    i.url as image_avatar
FROM
    users AS u
    LEFT JOIN images i ON u.img_id = i.id
;