CREATE OR REPLACE PROCEDURE delete_entity_group_content(_id integer, _entity_id integer, _entity_group_id integer)
    LANGUAGE plpgsql AS $$
DECLARE
    v_content_entity_id INT;
    v_content_entity_group_id INT;
BEGIN
    v_content_entity_id := _entity_id;
    v_content_entity_group_id := _entity_group_id;

    IF (_entity_id IS NULL AND _entity_group_id IS NULL) THEN
        SELECT content_entity_id, content_entity_group_id INTO v_content_entity_id, v_content_entity_group_id FROM entity_group_content WHERE id = _id;
    END IF;

    IF v_content_entity_id IS NOT NULL THEN
        DELETE FROM entities WHERE id = v_content_entity_id;
    END IF;

    IF v_content_entity_group_id IS NOT NULL THEN
        CALL delete_entity_group(v_content_entity_group_id);
    END IF;

    CALL delete_entity_group_content_and_move(_id);

END;
$$;