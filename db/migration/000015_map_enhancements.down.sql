
ALTER TABLE "map_pin_types" DROP COLUMN "section";
ALTER TABLE "map_pin_types" DROP COLUMN "map_pin_type_group_id";
ALTER TABLE "map_pin_types" ADD COLUMN map_id integer;

-- return data to the map_id column
UPDATE map_pin_types
    SET map_id = mp.map_id
FROM
    map_pin_types mpt
    JOIN map_pins mp ON mpt.id = mp.map_pin_type_id;

DO $$
DECLARE
    item RECORD;
BEGIN
    -- For each menu item
    FOR item IN (
        SELECT
            mp.id as map_pin_id,
            mp.map_id as map_pin_map_id,
            0 as new_map_pin_type_id,
            mpt.*
        FROM
            map_pins mp
            JOIN map_pin_types mpt ON mp.map_pin_type_id = mpt.id
        WHERE mp.map_id <> mpt.map_id
    ) LOOP

        -- create new map_pin_type with the same data, but correct map_id, then update map_pins to point to the new map_pin_type
        INSERT INTO map_pin_types (map_id, shape, background_color, border_color, icon_color, icon, icon_size, width
        ) VALUES (
                  item.map_pin_map_id,
                  item.shape,
                  item.background_color,
                  item.border_color,
                  item.icon_color,
                  item.icon,
                  item.icon_size,
                  item.width
        )
        RETURNING id INTO item.new_map_pin_type_id;

        UPDATE map_pins
            SET map_pin_type_id = item.new_map_pin_type_id
        WHERE id = item.map_pin_id;

    END LOOP;
END $$ LANGUAGE plpgsql;


ALTER TABLE "map_pin_types" ADD FOREIGN KEY ("map_id") REFERENCES "maps" ("id");

DROP TABLE world_map_pin_type_groups;
DROP TABLE map_pin_type_group;

DROP PROCEDURE IF EXISTS update_map_layer_is_main(INT);
