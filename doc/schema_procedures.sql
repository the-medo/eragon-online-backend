CREATE OR REPLACE PROCEDURE move_menu_item(p_id INT, p_target_position INT)
    LANGUAGE plpgsql AS $$
DECLARE
    v_old_position INT;
    v_menu_id INT;
    v_max_position INT;
BEGIN
    -- Get the current position and menu_id of the menu item
    SELECT "position", "menu_id" INTO v_old_position, v_menu_id
    FROM "menu_items"
    WHERE "id" = p_id;

    -- Get the maximum position within the menu
    SELECT MAX("position") INTO v_max_position
    FROM "menu_items"
    WHERE "menu_id" = v_menu_id;

    -- Check if the target position is valid
    IF p_target_position < 1 OR p_target_position > v_max_position THEN
        RAISE EXCEPTION 'Invalid target position';
    END IF;

    -- Update positions based on the move direction
    IF v_old_position < p_target_position THEN
        -- Move down
        UPDATE "menu_items"
        SET "position" = "position" - 1
        WHERE "menu_id" = v_menu_id
          AND "position" BETWEEN v_old_position + 1 AND p_target_position;

    ELSIF v_old_position > p_target_position THEN
        -- Move up
        UPDATE "menu_items"
        SET "position" = "position" + 1
        WHERE "menu_id" = v_menu_id
          AND "position" BETWEEN p_target_position AND v_old_position - 1;
    END IF;

    -- Set the new position of the menu item
    UPDATE "menu_items"
    SET "position" = p_target_position
    WHERE "id" = p_id;

END;
$$;

CREATE OR REPLACE PROCEDURE move_group_up(p_id INT)
    LANGUAGE plpgsql AS $$
DECLARE
    v_menu_id INT;
    v_target_group_start INT;
    v_target_group_end INT;
    v_prev_group_start INT;
    v_prev_group_end INT;
    v_target_group_size INT;
    v_prev_group_size INT;
    v_temp_offset INT := 1000;
BEGIN
    -- Get the menu_id and the position of the main item of the target group
    SELECT "menu_id", "position" INTO v_menu_id, v_target_group_start
    FROM "menu_items"
    WHERE "id" = p_id;

    -- Find the end position of the target group
    SELECT COALESCE(MIN("position"), v_target_group_start) - 1 INTO v_target_group_end
    FROM "menu_items"
    WHERE "menu_id" = v_menu_id
      AND "position" > v_target_group_start
      AND "is_main" = true;

    -- Find the end position of the previous group
    SELECT MAX("position") INTO v_prev_group_end
    FROM "menu_items"
    WHERE "menu_id" = v_menu_id
      AND "position" < v_target_group_start;

    -- If there's no previous group, exit the procedure
    IF v_prev_group_end IS NULL THEN
        RAISE NOTICE 'This group is already at the top';
        RETURN;
    END IF;

    -- Find the start position of the previous group
    SELECT COALESCE(MAX("position"), 1) INTO v_prev_group_start
    FROM "menu_items"
    WHERE "menu_id" = v_menu_id
      AND "position" < v_prev_group_end
      AND "is_main" = true;

    -- Calculate the size of both groups
    v_target_group_size := v_target_group_end - v_target_group_start + 1;
    v_prev_group_size := v_prev_group_end - v_prev_group_start + 1;

    -- Temporarily move the target group out of the way
    UPDATE "menu_items"
    SET "position" = "position" + v_temp_offset
    WHERE "menu_id" = v_menu_id
      AND "position" BETWEEN v_target_group_start AND v_target_group_end;

    -- Move the previous group down by the size of the target group
    UPDATE "menu_items"
    SET "position" = "position" + v_target_group_size
    WHERE "menu_id" = v_menu_id
      AND "position" BETWEEN v_prev_group_start AND v_prev_group_end;

    -- Move the target group up by the size of the previous group
    UPDATE "menu_items"
    SET "position" = "position" - v_prev_group_size - v_temp_offset
    WHERE "menu_id" = v_menu_id
      AND "position" BETWEEN v_target_group_start + v_temp_offset AND v_target_group_end + v_temp_offset;

END;
$$;