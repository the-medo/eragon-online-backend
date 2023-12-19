CREATE UNIQUE INDEX ON "entities" ("post_id", "module_id");
CREATE UNIQUE INDEX ON "entities" ("map_id", "module_id");
CREATE UNIQUE INDEX ON "entities" ("location_id", "module_id");
CREATE UNIQUE INDEX ON "entities" ("image_id", "module_id");

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
