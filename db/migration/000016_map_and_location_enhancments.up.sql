CREATE OR REPLACE PROCEDURE delete_location(p_location_id INT)
    LANGUAGE plpgsql
AS $$
DECLARE
    _egc RECORD;
BEGIN
    -- remove location from associated map pins
    UPDATE map_pins SET location_id = NULL WHERE location_id = p_location_id;

    -- Delete all the entities of the location
    FOR _egc IN (
        SELECT
            content.id as entity_group_content_id,
            e.id as entity_id
        FROM
            entity_group_content content
            JOIN entities e ON e.id = content.content_entity_id
        WHERE
            e.location_id = p_location_id
    ) LOOP
        CALL delete_entity_group_content(_egc.entity_group_content_id, _egc.entity_id, NULL);
    END LOOP;

    DELETE FROM world_locations WHERE location_id = p_location_id;

    DELETE FROM locations WHERE id = p_location_id;
END;
$$;

CREATE OR REPLACE PROCEDURE delete_map(p_map_id INT)
    LANGUAGE plpgsql
AS $$
DECLARE
    _egc RECORD;
BEGIN
    DELETE FROM map_pins WHERE map_id = p_map_id;
    DELETE FROM map_layers WHERE map_id = p_map_id;
    DELETE FROM map_layers WHERE map_id = p_map_id;

    -- Delete all the entities of the location
    FOR _egc IN (
        SELECT
            content.id as entity_group_content_id,
            e.id as entity_id
        FROM
            entity_group_content content
            JOIN entities e ON e.id = content.content_entity_id
        WHERE
            e.map_id = p_map_id
    ) LOOP
        CALL delete_entity_group_content(_egc.entity_group_content_id, _egc.entity_id, NULL);
    END LOOP;

    DELETE FROM world_maps WHERE map_id = p_map_id;

    DELETE FROM maps WHERE id = p_map_id;
END;
$$;

DROP VIEW view_locations;
CREATE VIEW view_locations AS
SELECT
    l.*,
    i.url as thumbnail_image_url,
    p.title as post_title
FROM
    locations l
    LEFT JOIN images i ON l.thumbnail_image_id = i.id
    LEFT JOIN posts p ON l.post_id = p.id
;

CREATE VIEW view_connections_world_posts AS
    SELECT
        wl.world_id as world_id,
        l.post_id as post_id,
        'locations' as helper_name,
        l.id as helper_id
    FROM
        world_locations wl
        JOIN locations l ON wl.location_id = l.id JOIN posts p ON l.post_id = p.id
    WHERE l.post_id IS NOT NULL

    UNION ALL

    SELECT
        wm.world_id as world_id,
        mip.post_id,
        'menu_item_posts',
        mip.menu_id
    FROM
        world_menu wm
        JOIN menu_item_posts mip ON mip.menu_id = wm.menu_id

    UNION ALL

    SELECT
        wm.world_id,
        mi.description_post_id,
        'menu_items',
        mi.id
    FROM
        world_menu wm
        JOIN menu_items mi ON mi.menu_id = wm.menu_id
    WHERE mi.description_post_id IS NOT NULL

    UNION ALL

    SELECT
        id,
        description_post_id,
        'worlds',
        id
    FROM
        worlds
    WHERE description_post_id IS NOT NULL

    UNION ALL

    SELECT
        wm.world_id,
        e.post_id,
        'entities',
        e.id
    FROM
        world_menu wm
        JOIN menu_items mi ON mi.menu_id = wm.menu_id
        JOIN get_recursive_entities(mi.entity_group_id) re ON 1 = 1
        JOIN entities e ON e.id = re.content_entity_id
    WHERE e.post_id IS NOT NULL
;

CREATE VIEW view_connections_menus AS
SELECT
    m.id as menu_id,
    COALESCE(wm.world_id, 0) as world_id,
    0 as quest_id,
    0 as character_id,
    0 as system_id
FROM
    menus m
    LEFT JOIN world_menu wm ON m.id = wm.menu_id
--     LEFT JOIN quest_menu q ON m.id = q.menu_id
--     LEFT JOIN character_menu c ON m.id = c.menu_id
--     LEFT JOIN system_menu s ON m.id = s.menu_id
;

CREATE TABLE "world_posts" (
    "world_id" int NOT NULL,
    "post_id" int NOT NULL
);
ALTER TABLE "world_posts" ADD FOREIGN KEY ("world_id") REFERENCES "worlds" ("id");
ALTER TABLE "world_posts" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id");
CREATE UNIQUE INDEX ON "world_posts" ("world_id", "post_id");

INSERT INTO world_posts (world_id, post_id)
SELECT world_id, post_id FROM view_connections_world_posts
ON CONFLICT (world_id, post_id) DO NOTHING
;

CREATE UNIQUE INDEX ON "world_map_pin_type_groups" ("world_id", "map_pin_type_group_id");
CREATE UNIQUE INDEX ON "world_maps" ("world_id", "map_id");
CREATE UNIQUE INDEX ON "world_locations" ("world_id", "location_id");

CREATE OR REPLACE PROCEDURE assign_post_by_menu_id(p_post_id integer, p_menu_id integer)
    LANGUAGE plpgsql AS $$
DECLARE
    v_world_id INT;
    v_quest_id INT;
    v_character_id INT;
    v_system_id INT;
BEGIN

    SELECT "world_id", "quest_id", "character_id", "system_id" INTO v_world_id, v_quest_id, v_character_id, v_system_id
    FROM "view_connections_menus"
    WHERE "menu_id" = p_menu_id;

    IF v_world_id > 0 THEN
        INSERT INTO world_posts (world_id, post_id) VALUES (v_world_id, p_post_id) ON CONFLICT DO NOTHING;
    END IF;

    IF v_quest_id > 0 THEN
        RAISE EXCEPTION 'Assigning posts into quest is not implemented yet';
    END IF;

    IF v_character_id > 0 THEN
        RAISE EXCEPTION 'Assigning posts into character is not implemented yet';
    END IF;

    IF v_system_id > 0 THEN
        RAISE EXCEPTION 'Assigning posts into system is not implemented yet';
    END IF;

END;
$$;