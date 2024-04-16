CREATE OR REPLACE PROCEDURE move_map_layer(p_id INT, p_target_position INT)
    LANGUAGE plpgsql AS $$
DECLARE
    v_old_position INT;
    v_map_id INT;
    v_max_position INT;
BEGIN
    -- Get the current position and menu_id of the menu item
    SELECT "position", "map_id" INTO v_old_position, v_map_id
    FROM "map_layers"
    WHERE "id" = p_id;

    -- Get the maximum position within the menu
    SELECT MAX("position") INTO v_max_position
    FROM "map_layers"
    WHERE "map_id" = v_map_id;

    -- Check if the target position is valid
    IF p_target_position < 1 OR p_target_position > v_max_position THEN
        RAISE EXCEPTION 'Invalid target position';
    END IF;

    -- Update positions based on the move direction
    IF v_old_position < p_target_position THEN
        -- Move down
        UPDATE "map_layers"
        SET "position" = "position" - 1
        WHERE "map_id" = v_map_id
          AND "position" BETWEEN v_old_position + 1 AND p_target_position;

    ELSIF v_old_position > p_target_position THEN
        -- Move up
        UPDATE "map_layers"
        SET "position" = "position" + 1
        WHERE "map_id" = v_map_id
          AND "position" BETWEEN p_target_position AND v_old_position - 1;
    END IF;

    -- Set the new position of the menu item
    UPDATE "map_layers"
    SET "position" = p_target_position
    WHERE "id" = p_id;

END;
$$;