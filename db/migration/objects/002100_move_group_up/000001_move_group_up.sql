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
    SELECT
        COALESCE(
                MIN("position"),
                (SELECT MAX(position) + 1 FROM menu_items WHERE menu_id = 2)
        ) - 1 INTO v_target_group_end
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
    SELECT COALESCE(MAX("position"), 0) INTO v_prev_group_start
    FROM "menu_items"
    WHERE "menu_id" = v_menu_id
      AND "position" <= v_prev_group_end
      AND "is_main" = true;


    -- If there's no group for previous item, exit the procedure
    IF v_prev_group_start = 0 THEN
        RAISE NOTICE 'Previous item has no group';
        RETURN;
    END IF;

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
