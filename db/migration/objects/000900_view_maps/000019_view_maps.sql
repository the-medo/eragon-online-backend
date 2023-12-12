CREATE VIEW view_maps AS
SELECT
    m.*,
    i.url as thumbnail_image_url,
    e.id as entity_id,
    e.module_id as module_id,
    e.module_type as module_type,
    e.module_type_id as module_type_id,
    e.tags as tags
FROM
    maps m
    LEFT JOIN images i ON m.thumbnail_image_id = i.id
    LEFT JOIN view_entities e ON e.map_id = m.id
;