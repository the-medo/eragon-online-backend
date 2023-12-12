CREATE VIEW view_maps AS
SELECT
    m.*,
    i.url as thumbnail_image_url
FROM
    maps m
    LEFT JOIN images i ON m.thumbnail_image_id = i.id
;