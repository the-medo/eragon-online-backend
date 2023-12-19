CREATE VIEW view_module_type_tags_available AS
SELECT
    mtta.*,
    cast(COUNT(mt.module_id) as integer) as count
FROM
    module_type_tags_available mtta
        LEFT JOIN module_tags mt ON mt.tag_id = mtta.id
GROUP BY
    mtta.id
;