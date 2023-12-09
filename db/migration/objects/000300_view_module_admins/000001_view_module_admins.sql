CREATE VIEW view_module_admins AS
SELECT
    m.*,
    ma.user_id,
    ma.approved,
    ma.super_admin,
    ma.allowed_entity_types,
    ma.allowed_menu
FROM
    modules m
        JOIN module_admins ma ON ma.module_id = m.id
;