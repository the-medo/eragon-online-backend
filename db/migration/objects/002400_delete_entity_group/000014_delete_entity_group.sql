CREATE OR REPLACE PROCEDURE delete_entity_group(_entity_group_id integer)
    LANGUAGE plpgsql AS $$
DECLARE
    egc RECORD;
    is_main_menu_item_group INT;
    entity_group_content_id INT;
BEGIN

    SELECT COUNT(*) FROM menu_items WHERE entity_group_id = _entity_group_id INTO is_main_menu_item_group;
    SELECT id FROM entity_group_content WHERE content_entity_group_id = _entity_group_id INTO entity_group_content_id;

    -- first, we delete contents of the entity group
    FOR egc IN (SELECT * FROM entity_group_content WHERE entity_group_id = _entity_group_id) LOOP
            CALL delete_entity_group_content(egc.id, egc.content_entity_id, egc.content_entity_group_id);
        END LOOP;

    IF entity_group_content_id IS NOT NULL THEN
        CALL delete_entity_group_content_and_move(entity_group_content_id);
    END IF;

    -- then, we delete the entity group itself if it is not a main menu item group
    IF is_main_menu_item_group IS NOT NULL AND is_main_menu_item_group = 0 THEN
        DELETE FROM entity_groups WHERE id = _entity_group_id;
    END IF;
END;
$$;