
DROP VIEW view_menu_item_posts;
DROP VIEW view_posts;

ALTER TABLE "posts" DROP COLUMN post_type_id;
ALTER TABLE "post_history" DROP COLUMN post_type_id;

DROP TABLE "post_types";

CREATE VIEW view_posts AS
SELECT
    p.*,
    i.url as thumbnail_img_url,
    e.id as entity_id,
    e.module_id as module_id,
    e.module_type as module_type,
    e.module_type_id as module_type_id,
    e.tags as tags
FROM
    posts p
    LEFT JOIN images i ON p.thumbnail_img_id = i.id
    LEFT JOIN view_entities e ON e.post_id = p.id
;

CREATE VIEW view_menu_item_posts AS
SELECT * FROM menu_item_posts mip JOIN view_posts vp ON mip.post_id = vp.id;