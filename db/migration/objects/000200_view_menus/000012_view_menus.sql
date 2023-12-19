CREATE VIEW view_menus AS
SELECT
    m.*,
    i.url as header_image_url
FROM
    menus m
    LEFT JOIN images i ON m.menu_header_img_id = i.id
;