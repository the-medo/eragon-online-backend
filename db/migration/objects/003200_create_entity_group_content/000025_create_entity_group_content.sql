CREATE OR REPLACE PROCEDURE create_entity_group_content(p_entity_group_id INT, p_content_entity_group_id INT, p_content_entity_id INT, p_new_position INT)
    LANGUAGE plpgsql AS $$
BEGIN
    -- Move up new entity group contents
    UPDATE "entity_group_content"
    SET "position" = "position" + 1
    WHERE "entity_group_id" = p_entity_group_id
      AND "position" >= p_new_position;

    INSERT INTO "entity_group_content" (entity_group_id, position, content_entity_id, content_entity_group_id)  VALUES
    (p_entity_group_id, p_new_position, p_content_entity_id, p_content_entity_group_id);
END;
$$;