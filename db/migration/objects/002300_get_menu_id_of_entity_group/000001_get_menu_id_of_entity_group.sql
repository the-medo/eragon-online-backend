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