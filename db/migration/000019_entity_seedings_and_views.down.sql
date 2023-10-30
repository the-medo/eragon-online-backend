-- ============ return view_maps
DROP VIEW view_maps;
CREATE VIEW view_maps AS
SELECT
    m.*,
    i.url as thumbnail_image_url
FROM
    maps m
    LEFT JOIN images i ON m.thumbnail_image_id = i.id
;

-- ============ drop view_images
DROP VIEW view_images;

-- ============ return view_posts
DROP VIEW view_menu_item_posts;
DROP VIEW view_posts;

CREATE VIEW view_posts AS
SELECT
    p.*,
    pt.name as post_type_name,
    pt.draftable as post_type_draftable,
    pt.privatable as post_type_privatable,
    i.url as thumbnail_img_url
FROM
    posts p
    JOIN post_types pt ON p.post_type_id = pt.id
    LEFT JOIN images i ON p.thumbnail_img_id = i.id
;

CREATE VIEW view_menu_item_posts AS
SELECT * FROM menu_item_posts mip JOIN view_posts vp ON mip.post_id = vp.id;


-- ============ return view_locations

DROP VIEW view_locations;
CREATE VIEW view_locations AS
SELECT
    l.*,
    i.url as thumbnail_image_url,
    p.title as post_title
FROM
    locations l
    LEFT JOIN images i ON l.thumbnail_image_id = i.id
    LEFT JOIN posts p ON l.post_id = p.id
;


-- ============ drop new view and indexes
DROP VIEW IF EXISTS view_entities;

drop index entities_map_id_module_id_idx;
drop index entities_location_id_module_id_idx;
drop index entities_post_id_module_id_idx;
drop index entities_image_id_module_id_idx;
