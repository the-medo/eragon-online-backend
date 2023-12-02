
--====================================

CREATE OR REPLACE FUNCTION get_recursive_entities(_main_entity_group_id INT)
    RETURNS TABLE (
                      id INT,
                      entity_group_id INT,
                      content_entity_id INT,
                      content_entity_group_id INT
                  ) AS $$
BEGIN
    RETURN QUERY
        WITH RECURSIVE entity_recursive AS (
            SELECT
                egc.id,
                egc.entity_group_id,
                egc.content_entity_id,
                egc.content_entity_group_id
            FROM
                entity_group_content egc
            WHERE
                    egc.entity_group_id = _main_entity_group_id

            UNION ALL

            SELECT
                child_egc.id,
                child_egc.entity_group_id,
                child_egc.content_entity_id,
                child_egc.content_entity_group_id
            FROM
                entity_recursive er
                    JOIN entity_group_content child_egc ON er.content_entity_group_id = child_egc.entity_group_id
            WHERE
                child_egc.content_entity_id IS NOT NULL OR child_egc.content_entity_group_id IS NOT NULL
        )
        SELECT * FROM entity_recursive;
END;
$$ LANGUAGE plpgsql;

--====================================

CREATE VIEW view_menus AS
SELECT
    m.*,
    i.url as header_image_url
FROM
    menus m
        LEFT JOIN images i ON m.menu_header_img_id = i.id
;

--====================================

CREATE VIEW view_module_admins AS
SELECT
    m.*,
    ma.user_id,
    ma.approved,
    ma.super_admin,
    ma.allowed_entity_types,
    ma.allowed_menu
FROM
    modules m
    JOIN module_admins ma ON ma.module_id = m.id
;

--====================================

CREATE VIEW view_module_type_tags_available AS
SELECT
    mtta.*,
    cast(COUNT(mt.module_id) as integer) as count
FROM
    module_type_tags_available mtta
    LEFT JOIN module_tags mt ON mt.tag_id = mtta.id
GROUP BY
    mtta.id
;

--====================================

CREATE VIEW view_users AS
SELECT
    u.*,
    i.id as avatar_image_id,
    i.url as avatar_image_url,
    i.img_guid as avatar_image_guid,
    p.deleted_at as introduction_post_deleted_at
FROM
    users AS u
    LEFT JOIN images i ON u.img_id = i.id
    LEFT JOIN posts p ON u.introduction_post_id = p.id
;

--====================================

CREATE VIEW view_entities AS
SELECT
    e.*,
    m.module_type as module_type,
    CAST(CASE m.module_type
             WHEN 'world' THEN m.world_id
             WHEN 'quest' THEN m.quest_id
             WHEN 'character' THEN m.character_id
             WHEN 'system' THEN m.system_id
        END as integer) as module_type_id,
    tags.tags as tags
FROM
    entities e
        LEFT JOIN modules m ON e.module_id = m.id
        LEFT JOIN (
        SELECT
            et.entity_id,
            cast(array_agg(tag_available.id) as int[]) AS tags
        FROM
            entity_tags et
                LEFT JOIN module_entity_tags_available tag_available ON tag_available.id = et.tag_id
        GROUP BY et.entity_id
    ) tags ON tags.entity_id = e.id
;

--====================================

CREATE VIEW view_locations AS
SELECT
    l.*,
    i.url as thumbnail_image_url,
    p.title as post_title,
    e.id as entity_id,
    e.module_id as module_id,
    e.module_type as module_type,
    e.module_type_id as module_type_id,
    e.tags as tags
FROM
    locations l
        JOIN view_entities e ON e.location_id = l.id
        LEFT JOIN images i ON l.thumbnail_image_id = i.id
        LEFT JOIN posts p ON l.post_id = p.id
;

--====================================

CREATE VIEW view_images AS
SELECT
    i.*,
    e.id as entity_id,
    e.module_id as module_id,
    e.module_type as module_type,
    e.module_type_id as module_type_id,
    e.tags as tags
FROM
    images i
        LEFT JOIN view_entities e ON e.image_id = i.id
;

--====================================

CREATE VIEW view_maps AS
SELECT
    m.*,
    i.url as thumbnail_image_url,
    e.id as entity_id,
    e.module_id as module_id,
    e.module_type as module_type,
    e.module_type_id as module_type_id,
    e.tags as tags
FROM
    maps m
        LEFT JOIN images i ON m.thumbnail_image_id = i.id
        LEFT JOIN view_entities e ON e.map_id = m.id
;

--====================================


CREATE VIEW view_posts AS
SELECT
    p.*,
    i.url as thumbnail_img_url,
    e.id as entity_id,
    e.module_id as module_id,
    e.module_type as module_type,
    e.module_type_id as module_type_id,
    e.tags as tags
FROM
    posts p
        LEFT JOIN images i ON p.thumbnail_img_id = i.id
        LEFT JOIN view_entities e ON e.post_id = p.id
;

--====================================
CREATE VIEW view_menu_item_posts AS
SELECT * FROM menu_item_posts mip JOIN view_posts vp ON mip.post_id = vp.id;

--====================================


CREATE VIEW view_map_layers AS
SELECT
    ml.*,
    i.url as image_url
FROM
    map_layers ml
        LEFT JOIN images i ON ml.image_id = i.id
;

--====================================

CREATE VIEW view_map_pins AS
SELECT
    mp.*,
    l.name as location_name,
    l.post_id as location_post_id,
    l.description as location_description,
    l.thumbnail_image_id as location_thumbnail_image_id,
    i.url as location_thumbnail_image_url
FROM
    map_pins mp
        LEFT JOIN locations l ON mp.location_id = l.id
        LEFT JOIN images i ON l.thumbnail_image_id = i.id
;

--====================================


CREATE VIEW view_connections_posts AS
SELECT
    e.module_id as module_id,
    l.post_id as post_id,
    'locations.id' as table_column,
    l.id as table_column_value
FROM
    entities e
        JOIN locations l ON l.post_id = e.post_id
WHERE l.post_id IS NOT NULL

UNION ALL

SELECT
    m.id,
    mi.description_post_id,
    'menu_items.description_post_id',
    mi.id
FROM
    modules m
        JOIN menu_items mi ON mi.menu_id = m.menu_id
WHERE mi.description_post_id IS NOT NULL

UNION ALL

SELECT
    m.id,
    description_post_id,
    'worlds.id',
    w.id
FROM
    worlds w
        JOIN modules m ON m.world_id = w.id AND m.module_type = 'world'
WHERE description_post_id IS NOT NULL

UNION ALL

SELECT
    m.id,
    e.post_id,
    'entities.id',
    e.id
FROM
    modules m
        JOIN menu_items mi ON mi.menu_id = m.menu_id
        JOIN get_recursive_entities(mi.entity_group_id) re ON 1 = 1
        JOIN entities e ON e.id = re.content_entity_id
WHERE e.post_id IS NOT NULL
;

--====================================

CREATE VIEW view_modules AS
SELECT m.id as id,
       m.world_id as world_id,
       m.system_id as system_id,
       m.character_id as character_id,
       m.quest_id as quest_id,
       m.module_type as module_type,
       m.menu_id as menu_id,
       m.header_img_id as header_img_id,
       m.thumbnail_img_id as thumbnail_img_id,
       m.avatar_img_id as avatar_img_id,
       cast(array_agg(tags.tag_id) as integer[]) AS tags
FROM
    modules m
    LEFT JOIN module_tags tags ON tags.module_id = m.id
GROUP BY m.id
;



--====================================

CREATE VIEW view_worlds AS
SELECT
    w.*,
    vm.id as module_id,
    vm.menu_id,
    vm.header_img_id,
    vm.thumbnail_img_id,
    vm.avatar_img_id,
    vm.tags
FROM
    worlds w
    JOIN view_modules vm ON w.id = vm.world_id
;

--====================================
CREATE OR REPLACE FUNCTION get_worlds(_is_public boolean, _tags integer[], _order_by varchar, _order_direction varchar, _limit int, _offset int)
    RETURNS SETOF view_worlds AS
$func$
DECLARE
    _limit_string VARCHAR;
BEGIN
    IF _order_by IS NULL THEN
        _order_by := 'created_at';
    END IF;

    IF _order_direction IS NULL OR (_order_direction <> 'ASC' AND _order_direction <> 'DESC') THEN
        _order_direction := 'DESC';
    END IF;

    _limit_string := '';
    IF _limit > 0 THEN
        _limit_string := 'LIMIT ' || _limit;
    END IF;

    RETURN QUERY EXECUTE format('
        SELECT * FROM view_worlds
        WHERE
            ($1 IS NULL OR public = $1) AND
            (array_length($2, 1) IS NULL OR tags @> $2)
        ORDER BY %I ' || _order_direction || ' ' || _limit_string || '
        OFFSET $3', _order_by)
        USING _is_public, _tags, _offset;
END
$func$  LANGUAGE plpgsql;

--====================================
CREATE OR REPLACE FUNCTION get_posts(_is_private boolean, _is_draft boolean, _tags integer[], _user_id integer, _module_id integer, _module_type module_type, _order_by varchar, _order_direction varchar, _limit int, _offset int)
    RETURNS SETOF view_posts AS
$func$
DECLARE
    _limit_string VARCHAR;
BEGIN
    IF _order_by IS NULL THEN
        _order_by := 'created_at';
    END IF;

    IF _order_direction IS NULL OR (_order_direction <> 'ASC' AND _order_direction <> 'DESC') THEN
        _order_direction := 'DESC';
    END IF;

    _limit_string := '';
    IF _limit > 0 THEN
        _limit_string := 'LIMIT ' || _limit;
    END IF;

    RETURN QUERY EXECUTE format('
        SELECT * FROM view_posts
        WHERE
            ($1 IS NULL OR is_private = $1) AND
            ($2 IS NULL OR is_draft = $2) AND
            ($5 IS NULL OR user_id = $5) AND
            ($6 IS NULL OR module_id = $6) AND
            ($7 IS NULL OR module_type = $7) AND
            deleted_at IS NULL AND
            (array_length($3, 1) IS NULL OR tags @> $3)
        ORDER BY %I ' || _order_direction || ' ' || _limit_string || '
        OFFSET $4 ', _order_by)
        USING _is_private, _is_draft, _tags, _offset, _user_id, _module_id, _module_type;
END
$func$  LANGUAGE plpgsql;

--====================================
CREATE OR REPLACE FUNCTION get_images(_tags integer[], _width integer, _height integer, _user_id integer, _module_id integer, _module_type module_type, _order_by varchar, _order_direction varchar, _limit int, _offset int)
    RETURNS SETOF view_images AS
$func$
DECLARE
    _limit_string VARCHAR;
BEGIN
    IF _order_by IS NULL THEN
        _order_by := 'created_at';
    END IF;

    IF _order_direction IS NULL OR (_order_direction <> 'ASC' AND _order_direction <> 'DESC') THEN
        _order_direction := 'DESC';
    END IF;

    _limit_string := '';
    IF _limit > 0 THEN
        _limit_string := 'LIMIT ' || _limit;
    END IF;

    RETURN QUERY EXECUTE format('
        SELECT * FROM view_images
        WHERE
            ($1 IS NULL OR width = $1) AND
            ($2 IS NULL OR height = $2) AND
            ($5 IS NULL OR user_id = $5) AND
            ($6 IS NULL OR module_id = $6) AND
            ($7 IS NULL OR module_type = $7) AND
            (array_length($3, 1) IS NULL OR tags @> $3)
        ORDER BY %I ' || _order_direction || ' ' || _limit_string || '
        OFFSET $4 ', _order_by)
        USING _width, _height, _tags, _offset, _user_id, _module_id, _module_type;
END
$func$  LANGUAGE plpgsql;

--====================================

CREATE OR REPLACE PROCEDURE update_map_layer_is_main(map_layer_id INT)
    LANGUAGE plpgsql
AS $$
DECLARE
    _map_id INT;
BEGIN
    -- Getting the map_id for the given map_layer_id
    SELECT map_id INTO _map_id FROM map_layers WHERE id = map_layer_id;

    -- If no map_id is found, then exit
    IF _map_id IS NULL THEN
        RAISE NOTICE 'No such map_layer_id exists: %', map_layer_id;
        RETURN;
    END IF;

    -- Set is_main to false for all rows related to the retrieved map_id
    UPDATE map_layers
    SET is_main = false
    WHERE map_id = _map_id;

    -- Set is_main to true for the row with the specified map_layer_id
    UPDATE map_layers
    SET is_main = true
    WHERE id = map_layer_id;
END;
$$;

--====================================


CREATE OR REPLACE PROCEDURE move_menu_item_post(mi_id INT, p_id INT, p_target_position INT)
    LANGUAGE plpgsql AS $$
DECLARE
    v_old_position INT;
    v_menu_item_id INT;
    v_max_position INT;
BEGIN
    -- Get the current position and menu_item_id of the menu item post
    SELECT "position", "menu_item_id" INTO v_old_position, v_menu_item_id
    FROM "menu_item_posts"
    WHERE "menu_item_id" = mi_id AND "post_id" = p_id;

    -- Get the maximum position within the menu
    SELECT MAX("position") INTO v_max_position
    FROM "menu_item_posts"
    WHERE "menu_item_id" = v_menu_item_id;

    -- Check if the target position is valid
    IF p_target_position < 1 OR p_target_position > v_max_position THEN
        RAISE EXCEPTION 'Invalid target position';
    END IF;

    -- Update positions based on the move direction
    IF v_old_position < p_target_position THEN
        -- Move down
        UPDATE "menu_item_posts"
        SET "position" = "position" - 1
        WHERE "menu_item_id" = v_menu_item_id
          AND "position" BETWEEN v_old_position + 1 AND p_target_position;

    ELSIF v_old_position > p_target_position THEN
        -- Move up
        UPDATE "menu_item_posts"
        SET "position" = "position" + 1
        WHERE "menu_item_id" = v_menu_item_id
          AND "position" BETWEEN p_target_position AND v_old_position - 1;
    END IF;

    -- Set the new position of the menu item
    UPDATE "menu_item_posts"
    SET "position" = p_target_position
    WHERE "menu_item_id" = mi_id AND "post_id" = p_id;

END;
$$;

--====================================

CREATE OR REPLACE PROCEDURE move_menu_item(p_id INT, p_target_position INT)
    LANGUAGE plpgsql AS $$
DECLARE
    v_old_position INT;
    v_menu_id INT;
    v_max_position INT;
BEGIN
    -- Get the current position and menu_id of the menu item
    SELECT "position", "menu_id" INTO v_old_position, v_menu_id
    FROM "menu_items"
    WHERE "id" = p_id;

    -- Get the maximum position within the menu
    SELECT MAX("position") INTO v_max_position
    FROM "menu_items"
    WHERE "menu_id" = v_menu_id;

    -- Check if the target position is valid
    IF p_target_position < 1 OR p_target_position > v_max_position THEN
        RAISE EXCEPTION 'Invalid target position';
    END IF;

    -- Update positions based on the move direction
    IF v_old_position < p_target_position THEN
        -- Move down
        UPDATE "menu_items"
        SET "position" = "position" - 1
        WHERE "menu_id" = v_menu_id
          AND "position" BETWEEN v_old_position + 1 AND p_target_position;

    ELSIF v_old_position > p_target_position THEN
        -- Move up
        UPDATE "menu_items"
        SET "position" = "position" + 1
        WHERE "menu_id" = v_menu_id
          AND "position" BETWEEN p_target_position AND v_old_position - 1;
    END IF;

    -- Set the new position of the menu item
    UPDATE "menu_items"
    SET "position" = p_target_position
    WHERE "id" = p_id;

END;
$$;

--====================================

CREATE OR REPLACE PROCEDURE move_group_up(p_id INT)
    LANGUAGE plpgsql AS $$
DECLARE
    v_menu_id INT;
    v_target_group_start INT;
    v_target_group_end INT;
    v_prev_group_start INT;
    v_prev_group_end INT;
    v_target_group_size INT;
    v_prev_group_size INT;
    v_temp_offset INT := 1000;
BEGIN
    -- Get the menu_id and the position of the main item of the target group
    SELECT "menu_id", "position" INTO v_menu_id, v_target_group_start
    FROM "menu_items"
    WHERE "id" = p_id;

    -- Find the end position of the target group
    SELECT
            COALESCE(
                    MIN("position"),
                    (SELECT MAX(position) + 1 FROM menu_items WHERE menu_id = 2)
            ) - 1 INTO v_target_group_end
    FROM "menu_items"
    WHERE "menu_id" = v_menu_id
      AND "position" > v_target_group_start
      AND "is_main" = true;

    -- Find the end position of the previous group
    SELECT MAX("position") INTO v_prev_group_end
    FROM "menu_items"
    WHERE "menu_id" = v_menu_id
      AND "position" < v_target_group_start;

    -- If there's no previous group, exit the procedure
    IF v_prev_group_end IS NULL THEN
        RAISE NOTICE 'This group is already at the top';
        RETURN;
    END IF;

    -- Find the start position of the previous group
    SELECT COALESCE(MAX("position"), 0) INTO v_prev_group_start
    FROM "menu_items"
    WHERE "menu_id" = v_menu_id
      AND "position" <= v_prev_group_end
      AND "is_main" = true;


    -- If there's no group for previous item, exit the procedure
    IF v_prev_group_start = 0 THEN
        RAISE NOTICE 'Previous item has no group';
        RETURN;
    END IF;

    -- Calculate the size of both groups
    v_target_group_size := v_target_group_end - v_target_group_start + 1;
    v_prev_group_size := v_prev_group_end - v_prev_group_start + 1;

    -- Temporarily move the target group out of the way
    UPDATE "menu_items"
    SET "position" = "position" + v_temp_offset
    WHERE "menu_id" = v_menu_id
      AND "position" BETWEEN v_target_group_start AND v_target_group_end;

    -- Move the previous group down by the size of the target group
    UPDATE "menu_items"
    SET "position" = "position" + v_target_group_size
    WHERE "menu_id" = v_menu_id
      AND "position" BETWEEN v_prev_group_start AND v_prev_group_end;

    -- Move the target group up by the size of the previous group
    UPDATE "menu_items"
    SET "position" = "position" - v_prev_group_size - v_temp_offset
    WHERE "menu_id" = v_menu_id
      AND "position" BETWEEN v_target_group_start + v_temp_offset AND v_target_group_end + v_temp_offset;

END;
$$;

--====================================

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

--====================================

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
                menu_items meg
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

--====================================

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

--====================================

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

--====================================

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

--====================================

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

--====================================

CREATE OR REPLACE PROCEDURE delete_location(p_location_id INT)
    LANGUAGE plpgsql
AS $$
DECLARE
    _egc RECORD;
BEGIN
    -- remove location from associated map pins
    UPDATE map_pins SET location_id = NULL WHERE location_id = p_location_id;

    -- Delete all the entities of the location
    FOR _egc IN (
        SELECT
            content.id as entity_group_content_id,
            e.id as entity_id
        FROM
            entity_group_content content
                JOIN entities e ON e.id = content.content_entity_id
        WHERE
                e.location_id = p_location_id
    ) LOOP
            CALL delete_entity_group_content(_egc.entity_group_content_id, _egc.entity_id, NULL);
        END LOOP;

    DELETE FROM world_locations WHERE location_id = p_location_id;

    DELETE FROM locations WHERE id = p_location_id;
END;
$$;

--====================================

CREATE OR REPLACE PROCEDURE delete_map(p_map_id INT)
    LANGUAGE plpgsql
AS $$
DECLARE
    _egc RECORD;
BEGIN
    DELETE FROM map_pins WHERE map_id = p_map_id;
    DELETE FROM map_layers WHERE map_id = p_map_id;
    DELETE FROM map_layers WHERE map_id = p_map_id;

    -- Delete all the entities of the location
    FOR _egc IN (
        SELECT
            content.id as entity_group_content_id,
            e.id as entity_id
        FROM
            entity_group_content content
                JOIN entities e ON e.id = content.content_entity_id
        WHERE
                e.map_id = p_map_id
    ) LOOP
            CALL delete_entity_group_content(_egc.entity_group_content_id, _egc.entity_id, NULL);
        END LOOP;

    DELETE FROM world_maps WHERE map_id = p_map_id;

    DELETE FROM maps WHERE id = p_map_id;
END;
$$;

--====================================





















