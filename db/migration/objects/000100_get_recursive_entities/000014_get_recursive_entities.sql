CREATE OR REPLACE FUNCTION get_recursive_entities(_main_entity_group_id INT)
    RETURNS TABLE (
                      "id" INT,
                      "entity_group_id" INT,
                      "content_entity_id" INT,
                      "content_entity_group_id" INT,
                      "position" INT
                  ) AS $$
BEGIN
    RETURN QUERY
        WITH RECURSIVE entity_recursive AS (
            SELECT
                egc.id,
                egc.entity_group_id,
                egc.content_entity_id,
                egc.content_entity_group_id,
                egc.position
            FROM
                entity_group_content egc
            WHERE
                egc.entity_group_id = _main_entity_group_id

            UNION ALL

            SELECT
                child_egc.id,
                child_egc.entity_group_id,
                child_egc.content_entity_id,
                child_egc.content_entity_group_id,
                child_egc.position
            FROM
                entity_recursive er
                JOIN entity_group_content child_egc ON er.content_entity_group_id = child_egc.entity_group_id
            WHERE
                child_egc.content_entity_id IS NOT NULL OR child_egc.content_entity_group_id IS NOT NULL
        )
        SELECT * FROM entity_recursive;
END;
$$ LANGUAGE plpgsql;
