CREATE VIEW view_locations AS
SELECT
    l.*,
    i.url as thumbnail_image_url
FROM
    locations l
    LEFT JOIN images i ON l.thumbnail_image_id = i.id
;