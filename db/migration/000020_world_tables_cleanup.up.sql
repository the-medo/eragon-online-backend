
CREATE TABLE "module_map_pin_type_groups" (
    "module_id" int NOT NULL,
    "map_pin_type_group_id" int NOT NULL
);

CREATE UNIQUE INDEX ON "module_map_pin_type_groups" ("module_id", "map_pin_type_group_id");
ALTER TABLE "module_map_pin_type_groups" ADD FOREIGN KEY ("module_id") REFERENCES "modules" ("id");
ALTER TABLE "module_map_pin_type_groups" ADD FOREIGN KEY ("map_pin_type_group_id") REFERENCES "map_pin_type_group" ("id");

INSERT INTO module_map_pin_type_groups (module_id, map_pin_type_group_id)
SELECT
    m.id as module_id,
    wmpg.map_pin_type_group_id as map_pin_type_group_id
FROM
    modules m
    JOIN world_map_pin_type_groups wmpg ON wmpg.world_id = m.world_id;



ALTER TABLE "modules"
    ADD COLUMN menu_id integer,
    ADD COLUMN header_img_id integer,
    ADD COLUMN thumbnail_img_id integer,
    ADD COLUMN avatar_img_id integer
;

ALTER TABLE "modules" ADD FOREIGN KEY ("menu_id") REFERENCES "menus" ("id");
ALTER TABLE "modules" ADD FOREIGN KEY ("header_img_id") REFERENCES "images" ("id");
ALTER TABLE "modules" ADD FOREIGN KEY ("thumbnail_img_id") REFERENCES "images" ("id");
ALTER TABLE "modules" ADD FOREIGN KEY ("avatar_img_id") REFERENCES "images" ("id");

WITH module_menus AS (
    SELECT
        m.id as module_id,
        wm.menu_id as menu_id,
        wi.header_img_id as header_img_id,
        wi.thumbnail_img_id as thumbnail_img_id,
        wi.avatar_img_id as avatar_img_id
    FROM
        modules m
        JOIN worlds w ON w.id = m.world_id
        JOIN world_menu wm ON wm.world_id = w.id
        JOIN world_images wi ON wi.world_id = w.id
)
UPDATE modules
    SET
        menu_id = module_menus.menu_id,
        header_img_id = module_menus.header_img_id,
        thumbnail_img_id = module_menus.thumbnail_img_id,
        avatar_img_id = module_menus.avatar_img_id
    FROM module_menus
    WHERE module_menus.module_id = modules.id
;

DROP VIEW view_connections_world_posts;
DROP VIEW view_connections_menus;
DROP FUNCTION get_worlds(boolean, integer[], varchar, varchar, int, int);
DROP VIEW view_worlds;

DROP TABLE world_menu;
DROP TABLE world_images;
DROP TABLE world_maps;
DROP TABLE world_locations;
DROP TABLE world_posts;
DROP TABLE world_map_pin_type_groups;
DROP TABLE world_activity;


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

CREATE VIEW view_modules AS
SELECT m.id as module_id,
       m.world_id as module_world_id,
       m.system_id as module_system_id,
       m.character_id as module_character_id,
       m.quest_id as module_quest_id,
       m.module_type as module_type,
       m.menu_id as menu_id,
       m.header_img_id as header_img_id,
       m.thumbnail_img_id as thumbnail_img_id,
       m.avatar_img_id as avatar_img_id,
       i_header.url as image_header,
       i_thumbnail.url as image_thumbnail,
       i_avatar.url as image_avatar,
       tags.tags AS tags
FROM
    modules m
        LEFT JOIN (
        SELECT
            mt.module_id,
            cast(array_agg(mt.tag_id) as integer[]) AS tags
        FROM
            module_tags mt
        GROUP BY mt.module_id
    ) tags ON tags.module_id = m.id
    LEFT JOIN images i_header on m.header_img_id = i_header.id
    LEFT JOIN images i_thumbnail on m.thumbnail_img_id = i_thumbnail.id
    LEFT JOIN images i_avatar on m.avatar_img_id = i_avatar.id
;

CREATE VIEW view_worlds AS
SELECT
    w.*,
    vm.*
FROM
    worlds w
    JOIN view_modules vm ON w.id = vm.module_world_id
;

CREATE OR REPLACE FUNCTION get_worlds(_is_public boolean, _tags integer[], _order_by varchar, _order_direction varchar, _limit int, _offset int)
    RETURNS SETOF view_worlds AS
$func$
BEGIN
    IF _order_by IS NULL THEN
        _order_by := 'created_at';
    END IF;

    IF _order_direction IS NULL OR (_order_direction <> 'ASC' AND _order_direction <> 'DESC') THEN
        _order_direction := 'DESC';
    END IF;

    RETURN QUERY EXECUTE format('
        SELECT * FROM view_worlds
        WHERE
            ($1 IS NULL OR public = $1) AND
            (array_length($2, 1) IS NULL OR tags @> $2)
        ORDER BY %I ' || _order_direction || '
        LIMIT $3
        OFFSET $4', _order_by)
        USING _is_public, _tags, _limit, _offset;
END
$func$  LANGUAGE plpgsql;