CREATE VIEW view_posts AS
SELECT
    p.*,
    pt.name as post_type_name,
    pt.draftable as post_type_draftable,
    pt.privatable as post_type_privatable
FROM
    posts p
    JOIN post_types pt ON p.post_type_id = pt.id
;