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