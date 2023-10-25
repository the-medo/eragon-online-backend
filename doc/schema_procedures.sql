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


CREATE OR REPLACE PROCEDURE move_menu_item_post(mi_id INT, p_id INT, p_target_position INT)
    LANGUAGE plpgsql AS $$
DECLARE
    v_old_position INT;
    v_menu_item_id INT;
    v_max_position INT;
BEGIN
    -- Get the current position and menu_item_id of the menu item post
    SELECT "position", "menu_item_id" INTO v_old_position, v_menu_item_id
    FROM "menu_item_posts"
    WHERE "menu_item_id" = mi_id AND "post_id" = p_id;

    -- Get the maximum position within the menu
    SELECT MAX("position") INTO v_max_position
    FROM "menu_item_posts"
    WHERE "menu_item_id" = v_menu_item_id;

    -- Check if the target position is valid
    IF p_target_position < 1 OR p_target_position > v_max_position THEN
        RAISE EXCEPTION 'Invalid target position';
    END IF;

    -- Update positions based on the move direction
    IF v_old_position < p_target_position THEN
        -- Move down
        UPDATE "menu_item_posts"
        SET "position" = "position" - 1
        WHERE "menu_item_id" = v_menu_item_id
          AND "position" BETWEEN v_old_position + 1 AND p_target_position;

    ELSIF v_old_position > p_target_position THEN
        -- Move up
        UPDATE "menu_item_posts"
        SET "position" = "position" + 1
        WHERE "menu_item_id" = v_menu_item_id
          AND "position" BETWEEN p_target_position AND v_old_position - 1;
    END IF;

    -- Set the new position of the menu item
    UPDATE "menu_item_posts"
    SET "position" = p_target_position
    WHERE "menu_item_id" = mi_id AND "post_id" = p_id;

END;
$$;

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