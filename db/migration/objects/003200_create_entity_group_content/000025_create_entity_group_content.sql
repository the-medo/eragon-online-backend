CREATE OR REPLACE PROCEDURE create_entity_group_content(p_entity_group_id INT, p_content_entity_group_id INT, p_content_entity_id INT, p_new_position INT)
    LANGUAGE plpgsql AS $$
DECLARE
    v_max_position INT;
    v_final_position INT;
BEGIN
    SELECT COALESCE(MAX("position"), 0) INTO v_max_position
    FROM "entity_group_content"
    WHERE "entity_group_id" = p_entity_group_id;

    -- If we didn't get new position as a parameter, we add it to the end
    IF p_new_position IS NULL THEN
        v_final_position = v_max_position + 1;
    ELSE
        -- Check if the target position is valid
        IF p_new_position < 1 OR p_new_position > v_max_position + 1 THEN
            RAISE EXCEPTION 'Invalid target position';
        END IF;

        v_final_position = p_new_position;

        -- Increase positions of entity group contents
        UPDATE "entity_group_content"
        SET "position" = "position" + 1
        WHERE "entity_group_id" = p_entity_group_id
          AND "position" >= v_final_position;
    END IF;

    INSERT INTO "entity_group_content" (entity_group_id, position, content_entity_id, content_entity_group_id)  VALUES
    (p_entity_group_id, v_final_position, p_content_entity_id, p_content_entity_group_id);

END
$$;