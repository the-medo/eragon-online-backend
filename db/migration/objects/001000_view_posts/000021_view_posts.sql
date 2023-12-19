CREATE VIEW view_posts AS
SELECT
    p.*,
    i.url as thumbnail_img_url,
    e.id as entity_id,
    e.module_id as module_id,
    e.module_type as module_type,
    e.module_type_id as module_type_id,
    e.tags as tags
FROM
    posts p
    LEFT JOIN images i ON p.thumbnail_img_id = i.id
    LEFT JOIN view_entities e ON e.post_id = p.id
;