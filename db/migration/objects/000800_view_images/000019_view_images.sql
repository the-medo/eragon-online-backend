CREATE VIEW view_images AS
SELECT
    i.*,
    e.id as entity_id,
    e.module_id as module_id,
    e.module_type as module_type,
    e.module_type_id as module_type_id,
    e.tags as tags
FROM
    images i
    LEFT JOIN view_entities e ON e.image_id = i.id
;