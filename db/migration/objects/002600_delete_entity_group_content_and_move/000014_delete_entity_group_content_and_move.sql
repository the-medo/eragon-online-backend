 CREATE OR REPLACE PROCEDURE delete_entity_group_content_and_move(_id integer)
    LANGUAGE plpgsql AS $$
BEGIN
    WITH deleted_entity_group_content AS (
        DELETE FROM "entity_group_content" d
            WHERE d.id = _id
            RETURNING *
    )
    UPDATE "entity_group_content"
    SET "position" = "position" - 1
    WHERE
        "entity_group_id" = (SELECT entity_group_id FROM deleted_entity_group_content)
      AND "position" > (SELECT position FROM deleted_entity_group_content);
END;
$$;