
CREATE OR REPLACE PROCEDURE move_menu_entity_groups(p_menu_id INT, p_entity_group_id INT, p_target_position INT)
    LANGUAGE plpgsql AS $$
DECLARE
    v_old_position INT;
    v_max_position INT;
BEGIN
    -- Get the current position and menu_id of the entity group
    SELECT "position" INTO v_old_position
    FROM "menu_item_entity_groups"
    WHERE "menu_id" = p_menu_id
      AND "entity_group_id" = p_entity_group_id;

    -- Get the maximum position within the menu
    SELECT MAX("position") INTO v_max_position
    FROM "menu_item_entity_groups"
    WHERE "menu_id" = p_menu_id;

    -- Check if the target position is valid
    IF p_target_position < 1 OR p_target_position > v_max_position THEN
        RAISE EXCEPTION 'Invalid target position';
    END IF;

    -- Update positions based on the move direction
    IF v_old_position < p_target_position THEN
        -- Move down
        UPDATE "menu_item_entity_groups"
        SET "position" = "position" - 1
        WHERE "menu_id" = p_menu_id
          AND "position" BETWEEN v_old_position + 1 AND p_target_position;

    ELSIF v_old_position > p_target_position THEN
        -- Move up
        UPDATE "menu_item_entity_groups"
        SET "position" = "position" + 1
        WHERE "menu_id" = p_menu_id
          AND "position" BETWEEN p_target_position AND v_old_position - 1;
    END IF;

    -- Set the new position of the entity group
    UPDATE "menu_item_entity_groups"
    SET "position" = p_target_position
    WHERE "menu_id" = p_menu_id
      AND "entity_group_id" = p_entity_group_id;

END;
$$;


-- content


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


create function get_menu_id_of_entity_group(_entity_group_id integer)
    returns TABLE(menu_id integer, entity_group_id integer)
    language plpgsql
as
$$
BEGIN
    RETURN QUERY
        WITH RECURSIVE entity_group_hierarchy AS (
            SELECT
                meg.menu_id,
                meg.entity_group_id
            FROM
                menu_item_entity_groups meg
            UNION ALL
            SELECT
                egh.menu_id,
                egc.content_entity_group_id
            FROM
                entity_group_hierarchy egh
                    JOIN entity_group_content egc ON egh.entity_group_id = egc.entity_group_id
            WHERE
                egc.content_entity_group_id IS NOT NULL
        )
        SELECT * FROM entity_group_hierarchy egh2 WHERE egh2.entity_group_id = _entity_group_id;
END;
$$;

alter function get_menu_id_of_entity_group(integer) owner to root;


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

CREATE OR REPLACE PROCEDURE delete_entity_group(_entity_group_id integer)
    LANGUAGE plpgsql AS $$
DECLARE
    egc RECORD;
    is_main_menu_item_group INT;
    entity_group_content_id INT;
BEGIN

    SELECT COUNT(*) FROM menu_item_entity_groups WHERE entity_group_id = _entity_group_id INTO is_main_menu_item_group;
    SELECT id FROM entity_group_content WHERE content_entity_group_id = _entity_group_id INTO entity_group_content_id;

    -- first, we delete contents of the entity group
    FOR egc IN (SELECT * FROM entity_group_content WHERE entity_group_id = _entity_group_id) LOOP
            CALL delete_entity_group_content(egc.id, egc.content_entity_id, egc.content_entity_group_id);
        END LOOP;

    IF entity_group_content_id > 0 THEN
        CALL delete_entity_group_content_and_move(entity_group_content_id);
    END IF;

    -- then, we delete the entity group itself if it is not a main menu item group
    IF is_main_menu_item_group == 0 THEN
        DELETE FROM entity_groups WHERE id = _entity_group_id;
    END IF;


END;
$$;