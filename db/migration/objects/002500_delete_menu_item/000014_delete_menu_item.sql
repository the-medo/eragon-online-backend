CREATE OR REPLACE PROCEDURE delete_menu_item(_menu_item_id INT)
    LANGUAGE plpgsql AS $$
DECLARE
    _entity_group_id INT;
    _menu_id INT;
    _position INT;
BEGIN
    -- Delete the menu item and get the entity_group_id
    WITH deleted_item AS (
        DELETE FROM "menu_items"
            WHERE id = _menu_item_id
            RETURNING entity_group_id, menu_id, position
    )
    SELECT entity_group_id, menu_id, position INTO _entity_group_id, _menu_id, _position FROM deleted_item;

    UPDATE "menu_items"
    SET "position" = "position" - 1
    WHERE "menu_id" = _menu_id AND "position" > _position;


    IF _entity_group_id IS NOT NULL THEN
        CALL delete_entity_group(_entity_group_id);
    END IF;
END;
$$;