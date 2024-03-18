CREATE OR REPLACE PROCEDURE delete_entity_group(_entity_group_id integer, _delete_type delete_entity_group_content_action)
    LANGUAGE plpgsql AS $$
DECLARE
    egc RECORD;
    v_is_main_menu_item_group INT;
    v_entity_group_content_id INT;
    v_parent_entity_group_id INT;
    v_position_in_parent INT;
    v_children_count INT;
BEGIN

    SELECT COUNT(*) FROM menu_items WHERE entity_group_id = _entity_group_id INTO v_is_main_menu_item_group;
    SELECT id, entity_group_id, position FROM entity_group_content WHERE content_entity_group_id = _entity_group_id INTO v_entity_group_content_id, v_parent_entity_group_id, v_position_in_parent;

    SELECT COUNT(*) FROM entity_group_content WHERE entity_group_id = _entity_group_id INTO v_children_count;

    IF _delete_type = 'delete_children' THEN

        -- first, we delete contents of the entity group
        FOR egc IN (SELECT * FROM entity_group_content WHERE entity_group_id = _entity_group_id) LOOP
                CALL delete_entity_group_content(egc.id, egc.content_entity_id, egc.content_entity_group_id, _delete_type);
            END LOOP;


        IF v_entity_group_content_id IS NOT NULL THEN
            CALL delete_entity_group_content_and_move(v_entity_group_content_id);
        END IF;

    end if;

    IF _delete_type = 'move_children' THEN

        UPDATE entity_group_content SET position = position + v_children_count - 1 WHERE entity_group_id = v_parent_entity_group_id AND position > v_position_in_parent;

        UPDATE entity_group_content SET position = position + v_position_in_parent - 1, entity_group_id = v_parent_entity_group_id WHERE entity_group_id = _entity_group_id;


        DELETE FROM "entity_group_content" d WHERE d.id = v_entity_group_content_id;

    end if;

    -- then, we delete the entity group itself if it is not a main menu item group
    IF v_is_main_menu_item_group IS NOT NULL AND v_is_main_menu_item_group = 0 THEN
        DELETE FROM entity_groups WHERE id = _entity_group_id;
    END IF;
END;
$$;