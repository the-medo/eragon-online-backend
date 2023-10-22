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
