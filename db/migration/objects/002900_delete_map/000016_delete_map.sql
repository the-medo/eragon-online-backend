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

    DELETE FROM entities WHERE map_id = p_map_id;

    DELETE FROM maps WHERE id = p_map_id;
END;
$$;