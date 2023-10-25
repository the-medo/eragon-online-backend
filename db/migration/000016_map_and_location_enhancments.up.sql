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


CREATE TABLE "world_posts" (
    "world_id" int NOT NULL,
    "post_id" int NOT NULL
);
ALTER TABLE "world_posts" ADD FOREIGN KEY ("world_id") REFERENCES "worlds" ("id");
ALTER TABLE "world_posts" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id");
CREATE UNIQUE INDEX ON "world_posts" ("world_id", "post_id");

CREATE UNIQUE INDEX ON "world_map_pin_type_groups" ("world_id", "map_pin_type_group_id");
CREATE UNIQUE INDEX ON "world_maps" ("world_id", "map_id");
CREATE UNIQUE INDEX ON "world_locations" ("world_id", "location_id");
