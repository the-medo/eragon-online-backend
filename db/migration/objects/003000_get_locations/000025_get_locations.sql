CREATE OR REPLACE FUNCTION get_locations(_tags integer[], _module_id integer, _module_type module_type, _order_by varchar, _order_direction varchar, _limit int, _offset int)
    RETURNS SETOF view_locations AS
$func$
DECLARE
    _limit_string VARCHAR;
BEGIN
    IF _order_by IS NULL THEN
        _order_by := 'id';
    END IF;

    IF _order_direction IS NULL OR (_order_direction <> 'ASC' AND _order_direction <> 'DESC') THEN
        _order_direction := 'DESC';
    END IF;

    _limit_string := '';
    IF _limit > 0 THEN
        _limit_string := 'LIMIT ' || _limit;
    END IF;

    RETURN QUERY EXECUTE format('
        SELECT * FROM view_locations
        WHERE
            ($3 IS NULL OR module_id = $3) AND
            ($4 IS NULL OR module_type = $4) AND
            (array_length($1, 1) IS NULL OR tags @> $1)
        ORDER BY %I ' || _order_direction || ' ' || _limit_string || '
        OFFSET $2 ', _order_by)
        USING _tags, _offset, _module_id, _module_type;
END
$func$  LANGUAGE plpgsql;