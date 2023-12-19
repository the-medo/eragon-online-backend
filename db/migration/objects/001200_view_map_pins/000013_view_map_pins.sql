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