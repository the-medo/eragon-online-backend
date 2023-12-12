CREATE VIEW view_map_layers AS
SELECT
    ml.*,
    i.url as image_url
FROM
    map_layers ml
    LEFT JOIN images i ON ml.image_id = i.id
;