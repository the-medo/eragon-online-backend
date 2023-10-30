CREATE UNIQUE INDEX ON "entities" ("post_id", "module_id");
CREATE UNIQUE INDEX ON "entities" ("map_id", "module_id");
CREATE UNIQUE INDEX ON "entities" ("location_id", "module_id");
CREATE UNIQUE INDEX ON "entities" ("image_id", "module_id");


DROP VIEW IF EXISTS view_entities;
CREATE VIEW view_entities AS
SELECT
    e.*,
    m.module_type as module_type,
    CAST(CASE m.module_type
        WHEN 'world' THEN m.world_id
        WHEN 'quest' THEN m.quest_id
        WHEN 'character' THEN m.character_id
        WHEN 'system' THEN m.system_id
    END as integer) as module_type_id,
    tags.tags as tags
FROM
    entities e
    LEFT JOIN modules m ON e.module_id = m.id
    LEFT JOIN (
        SELECT
            et.entity_id,
            cast(array_agg(tag_available.id) as int[]) AS tags
        FROM
            entity_tags et
            LEFT JOIN module_entity_tags_available tag_available ON tag_available.id = et.tag_id
        GROUP BY et.entity_id
    ) tags ON tags.entity_id = e.id
;

--  ============= entities table seeding ==================
--  ============= locations ==================

INSERT INTO entities (type, location_id, module_id)
SELECT
    'location' as type,
    l.id as location_id,
    m.id as module_id
FROM
    locations l
    JOIN world_locations wl ON wl.location_id = l.id
    JOIN worlds w ON w.id = wl.world_id
    JOIN modules m ON m.module_type = 'world' AND m.world_id = w.id
ON CONFLICT (location_id, module_id) DO NOTHING
;

--  ============= posts ==================
INSERT INTO entities (type, post_id, module_id)
SELECT
    'post' as type,
    p.id as post_id,
    m.id as module_id
FROM
    posts p
    JOIN world_posts wp ON wp.post_id = p.id
    JOIN worlds w ON w.id = wp.world_id
    JOIN modules m ON m.module_type = 'world' AND m.world_id = w.id
ON CONFLICT (post_id, module_id) DO NOTHING
;

--  ============= images from world_images ==================
INSERT INTO entities (type, image_id, module_id)
SELECT
    'image' as type,
    i.id    as image_id,
    m.id    as module_id
FROM
    images i
    JOIN world_images wi ON (wi.header_img_id = i.id OR wi.thumbnail_img_id = i.id OR wi.avatar_img_id = i.id)
    JOIN worlds w ON w.id = wi.world_id
    JOIN modules m ON m.module_type = 'world' AND m.world_id = w.id
ON CONFLICT (image_id, module_id) DO NOTHING
;

--  ============= images from locations ==================
INSERT INTO entities (type, image_id, module_id)
SELECT
    'image' as type,
    i.id    as image_id,
    m.id    as module_id
FROM
    images i
    JOIN locations l ON l.thumbnail_image_id = i.id
    JOIN world_locations wl ON wl.location_id = l.id
    JOIN worlds w ON w.id = wl.world_id
    JOIN modules m ON m.module_type = 'world' AND m.world_id = w.id
ON CONFLICT (image_id, module_id) DO NOTHING
;

--  ============= images from post thumbnails ==================
INSERT INTO entities (type, image_id, module_id)
SELECT
    'image' as type,
    i.id    as image_id,
    m.id    as module_id
FROM
    images i
    JOIN posts p ON p.thumbnail_img_id = i.id
    JOIN world_posts wp ON wp.post_id = p.id
    JOIN worlds w ON w.id = wp.world_id
    JOIN modules m ON m.module_type = 'world' AND m.world_id = w.id
ON CONFLICT (image_id, module_id) DO NOTHING
;

-- ============ rework view_locations

DROP VIEW view_locations;
CREATE VIEW view_locations AS
SELECT
    l.*,
    i.url as thumbnail_image_url,
    p.title as post_title,
    e.id as entity_id,
    e.module_id as module_id,
    e.module_type as module_type,
    e.module_type_id as module_type_id,
    e.tags as tags
FROM
    locations l
    JOIN view_entities e ON e.location_id = l.id
    LEFT JOIN images i ON l.thumbnail_image_id = i.id
    LEFT JOIN posts p ON l.post_id = p.id
;

-- ============ rework view_posts

DROP VIEW view_menu_item_posts;
DROP VIEW view_posts;

CREATE VIEW view_posts AS
SELECT
    p.*,
    pt.name as post_type_name,
    pt.draftable as post_type_draftable,
    pt.privatable as post_type_privatable,
    i.url as thumbnail_img_url,
    e.id as entity_id,
    e.module_id as module_id,
    e.module_type as module_type,
    e.module_type_id as module_type_id,
    e.tags as tags
FROM
    posts p
    JOIN post_types pt ON p.post_type_id = pt.id
    LEFT JOIN images i ON p.thumbnail_img_id = i.id
    LEFT JOIN view_entities e ON e.post_id = p.id
;

CREATE VIEW view_menu_item_posts AS
SELECT * FROM menu_item_posts mip JOIN view_posts vp ON mip.post_id = vp.id;


-- ============ add view_images
CREATE VIEW view_images AS
SELECT
    i.*,
    e.id as entity_id,
    e.module_id as module_id,
    e.module_type as module_type,
    e.module_type_id as module_type_id,
    e.tags as tags
FROM
    images i
    LEFT JOIN view_entities e ON e.image_id = i.id
;


-- ============ rework view_maps
DROP VIEW view_maps;
CREATE VIEW view_maps AS
SELECT
    m.*,
    i.url as thumbnail_image_url,
    e.id as entity_id,
    e.module_id as module_id,
    e.module_type as module_type,
    e.module_type_id as module_type_id,
    e.tags as tags
FROM
    maps m
    LEFT JOIN images i ON m.thumbnail_image_id = i.id
    LEFT JOIN view_entities e ON e.map_id = m.id
;