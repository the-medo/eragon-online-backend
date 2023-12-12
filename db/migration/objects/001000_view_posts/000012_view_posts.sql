CREATE VIEW view_posts AS
SELECT
    p.*,
    pt.name as post_type_name,
    pt.draftable as post_type_draftable,
    pt.privatable as post_type_privatable,
    i.url as thumbnail_img_url
FROM
    posts p
    JOIN post_types pt ON p.post_type_id = pt.id
    LEFT JOIN images i ON p.thumbnail_img_id = i.id
;