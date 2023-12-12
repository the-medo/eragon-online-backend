CREATE OR REPLACE PROCEDURE update_map_layer_is_main(map_layer_id INT)
    LANGUAGE plpgsql
AS $$
DECLARE
    _map_id INT;
BEGIN
    -- Getting the map_id for the given map_layer_id
    SELECT map_id INTO _map_id FROM map_layers WHERE id = map_layer_id;

    -- If no map_id is found, then exit
    IF _map_id IS NULL THEN
        RAISE NOTICE 'No such map_layer_id exists: %', map_layer_id;
        RETURN;
    END IF;

    -- Set is_main to false for all rows related to the retrieved map_id
    UPDATE map_layers
    SET is_main = false
    WHERE map_id = _map_id;

    -- Set is_main to true for the row with the specified map_layer_id
    UPDATE map_layers
    SET is_main = true
    WHERE id = map_layer_id;
END;
$$;