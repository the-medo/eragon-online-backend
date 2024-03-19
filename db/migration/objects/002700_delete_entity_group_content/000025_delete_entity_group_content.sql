CREATE OR REPLACE PROCEDURE delete_entity_group_content(_id integer, _delete_type delete_entity_group_content_action)
    LANGUAGE plpgsql AS $$
DECLARE
    v_content_entity_id INT;
    v_content_entity_group_id INT;
BEGIN
    SELECT content_entity_id, content_entity_group_id INTO v_content_entity_id, v_content_entity_group_id FROM entity_group_content WHERE id = _id;

    IF v_content_entity_id IS NOT NULL THEN
--         DELETE FROM entities WHERE id = v_content_entity_id;
        CALL delete_entity_group_content_and_move(_id);
    END IF;

    IF v_content_entity_group_id IS NOT NULL THEN
        CALL delete_entity_group(v_content_entity_group_id, _delete_type);
    END IF;

END;
$$;