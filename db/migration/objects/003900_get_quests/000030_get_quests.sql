CREATE OR REPLACE FUNCTION get_quests(
    _is_public boolean,
    _tags integer[],
    _world_id int,
    _system_id int,
    _order_by varchar,
    _order_direction varchar,
    _limit int,
    _offset int
) RETURNS SETOF view_quests AS
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
        SELECT * FROM view_quests
        WHERE
            ($1 IS NULL OR public = $1) AND
            (array_length($2, 1) IS NULL OR tags @> $2) AND
            ($3 IS NULL OR world_id = $3) AND
            ($4 IS NULL OR system_id = $4)
        ORDER BY %I %s %s
        OFFSET $5', _order_by, _order_direction, _limit_string)
    USING _is_public, _tags, _world_id, _system_id, _offset;
END
$func$  LANGUAGE plpgsql;