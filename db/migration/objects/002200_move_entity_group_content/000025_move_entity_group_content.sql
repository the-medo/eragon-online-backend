CREATE OR REPLACE PROCEDURE move_entity_group_content(p_id INT, p_new_entity_group_id INT, p_new_position INT)
    LANGUAGE plpgsql AS $$
DECLARE
    v_old_position INT;
    v_old_entity_group_id INT;
    v_new_max_position INT;
BEGIN
    -- Get the current position and entity_group_id of the content
    SELECT "position", "entity_group_id" INTO v_old_position, v_old_entity_group_id
    FROM "entity_group_content"
    WHERE "id" = p_id;

    -- in case of empty entity group, we move inside of the same group
    IF p_new_entity_group_id IS NULL THEN
        p_new_entity_group_id := v_old_entity_group_id;
    end if;

    IF p_new_entity_group_id != v_old_entity_group_id THEN

        -- Get the maximum position within the new group
        SELECT MAX("position") INTO v_new_max_position
        FROM "entity_group_content"
        WHERE "entity_group_id" = p_new_entity_group_id;

        IF p_new_position IS NULL THEN
            p_new_position := v_new_max_position;
        end if;

        -- Check if the target position is valid
        IF p_new_position < 1 OR p_new_position > v_new_max_position THEN
            RAISE EXCEPTION 'Invalid target position';
        END IF;

        -- Move down old entity group contents
        UPDATE "entity_group_content"
        SET "position" = "position" - 1
        WHERE "entity_group_id" = v_old_entity_group_id
          AND "position" >= v_old_position;

        -- Move up new entity group contents
        UPDATE "entity_group_content"
        SET "position" = "position" + 1
        WHERE "entity_group_id" = p_new_entity_group_id
          AND "position" >= p_new_position;

    ELSE -- group is moved inside of the same group

        -- Move up new entity group contents
        UPDATE "entity_group_content"
        SET "position" = "position" + CASE WHEN p_new_position > v_old_position THEN -1 ELSE 1 END
        WHERE "entity_group_id" = v_old_entity_group_id
          AND "position" BETWEEN p_new_position AND v_old_position;

    END IF;

    -- Set the new position and parent of the group content
    UPDATE "entity_group_content"
    SET "position" = p_new_position, "entity_group_id" = p_new_entity_group_id
    WHERE "id" = p_id;

END;
$$;