CREATE OR REPLACE PROCEDURE move_entity_group_content(p_id INT, p_target_position INT)
    LANGUAGE plpgsql AS $$
DECLARE
    v_old_position INT;
    v_entity_group_id INT;
    v_max_position INT;
BEGIN
    -- Get the current position and entity_group_id of the content
    SELECT "position", "entity_group_id" INTO v_old_position, v_entity_group_id
    FROM "entity_group_content"
    WHERE "id" = p_id;

    -- Get the maximum position within the group
    SELECT MAX("position") INTO v_max_position
    FROM "entity_group_content"
    WHERE "entity_group_id" = v_entity_group_id;

    -- Check if the target position is valid
    IF p_target_position < 1 OR p_target_position > v_max_position THEN
        RAISE EXCEPTION 'Invalid target position';
    END IF;

    -- Update positions based on the move direction
    IF v_old_position < p_target_position THEN
        -- Move down
        UPDATE "entity_group_content"
        SET "position" = "position" - 1
        WHERE "entity_group_id" = v_entity_group_id
          AND "position" BETWEEN v_old_position + 1 AND p_target_position;

    ELSIF v_old_position > p_target_position THEN
        -- Move up
        UPDATE "entity_group_content"
        SET "position" = "position" + 1
        WHERE "entity_group_id" = v_entity_group_id
          AND "position" BETWEEN p_target_position AND v_old_position - 1;
    END IF;

    -- Set the new position of the group content
    UPDATE "entity_group_content"
    SET "position" = p_target_position
    WHERE "id" = p_id;

END;
$$;