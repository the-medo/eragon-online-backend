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