
CREATE VIEW view_connections_posts AS
SELECT
    e.module_id as module_id,
    l.post_id as post_id,
    'locations.id' as table_column,
    l.id as table_column_value
FROM
    entities e
    JOIN locations l ON l.post_id = e.post_id
WHERE l.post_id IS NOT NULL

UNION ALL

SELECT
    m.id,
    mi.description_post_id,
    'menu_items.description_post_id',
    mi.id
FROM
    modules m
    JOIN menu_items mi ON mi.menu_id = m.menu_id
WHERE mi.description_post_id IS NOT NULL

UNION ALL

SELECT
    m.id,
    description_post_id,
    'worlds.id',
    w.id
FROM
    worlds w
    JOIN modules m ON m.world_id = w.id AND m.module_type = 'world'
WHERE description_post_id IS NOT NULL

UNION ALL

SELECT
    m.id,
    e.post_id,
    'entities.id',
    e.id
FROM
    modules m
    JOIN menu_items mi ON mi.menu_id = m.menu_id
    JOIN get_recursive_entities(mi.entity_group_id) re ON 1 = 1
    JOIN entities e ON e.id = re.content_entity_id
WHERE e.post_id IS NOT NULL
;