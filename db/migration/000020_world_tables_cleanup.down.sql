DROP VIEW view_module_admins;

CREATE TABLE "world_admins" (
    "world_id" int NOT NULL,
    "user_id" int NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "super_admin" boolean NOT NULL DEFAULT false,
    "approved" int NOT NULL,
    "motivational_letter" varchar NOT NULL
);

CREATE UNIQUE INDEX ON "world_admins" ("world_id", "user_id");
COMMENT ON COLUMN "world_admins"."approved" IS '0 = NO, 1 = YES, 2 = PENDING';
ALTER TABLE "world_admins" ADD FOREIGN KEY ("world_id") REFERENCES "worlds" ("id");
ALTER TABLE "world_admins" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

INSERT INTO world_admins (world_id, user_id, super_admin, approved, motivational_letter)
SELECT
    m.id as module_id,
    ma.user_id as user_id,
    ma.super_admin as super_admin,
    ma.approved as approved,
    ma.motivational_letter as motivational_letter
FROM
    modules m
    JOIN module_admins ma ON ma.module_id = m.world_id
WHERE
    m.module_type = 'world';

DROP TABLE module_admins;


CREATE TABLE "world_map_pin_type_groups" (
    "world_id" int NOT NULL,
    "map_pin_type_group_id" int NOT NULL
);
CREATE UNIQUE INDEX ON "world_map_pin_type_groups" ("world_id", "map_pin_type_group_id");
ALTER TABLE "world_map_pin_type_groups" ADD FOREIGN KEY ("world_id") REFERENCES "worlds" ("id");
ALTER TABLE "world_map_pin_type_groups" ADD FOREIGN KEY ("map_pin_type_group_id") REFERENCES "map_pin_type_group" ("id");

INSERT INTO world_map_pin_type_groups (world_id, map_pin_type_group_id)
SELECT
    m.world_id as world_id,
    mmpg.map_pin_type_group_id as map_pin_type_group_id
FROM
    modules m
    JOIN module_map_pin_type_groups mmpg ON mmpg.module_id = m.id
WHERE
    m.module_type = 'world';

DROP TABLE module_map_pin_type_groups;


CREATE TABLE "world_activity" (
    "world_id" int NOT NULL,
    "date" date NOT NULL,
    "post_count" int NOT NULL,
    "quest_count" int NOT NULL,
    "resource_count" int NOT NULL
);
CREATE UNIQUE INDEX ON "world_activity" ("world_id", "date");
ALTER TABLE "world_activity" ADD FOREIGN KEY ("world_id") REFERENCES "worlds" ("id");


CREATE TABLE "world_posts" (
    "world_id" int NOT NULL,
    "post_id" int NOT NULL
);
CREATE UNIQUE INDEX ON "world_posts" ("world_id", "post_id");
ALTER TABLE "world_posts" ADD FOREIGN KEY ("world_id") REFERENCES "worlds" ("id");
ALTER TABLE "world_posts" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id");

INSERT INTO world_posts (world_id, post_id)
SELECT
    module_type_id as world_id,
    post_id as post_id
FROM
    view_entities
WHERE
    type='post' AND
    module_type = 'world'
;


CREATE TABLE "world_locations" (
    "world_id" int NOT NULL,
    "location_id" int NOT NULL
);
CREATE UNIQUE INDEX ON "world_locations" ("world_id", "location_id");
ALTER TABLE "world_locations" ADD FOREIGN KEY ("world_id") REFERENCES "worlds" ("id");
ALTER TABLE "world_locations" ADD FOREIGN KEY ("location_id") REFERENCES "locations" ("id");

INSERT INTO world_locations (world_id, location_id)
SELECT
    module_type_id as world_id,
    location_id as location_id
FROM
    view_entities
WHERE
    type='location' AND
    module_type = 'world'
;



CREATE TABLE "world_maps" (
    "world_id" int NOT NULL,
    "map_id" int NOT NULL
);
CREATE UNIQUE INDEX ON "world_maps" ("world_id", "map_id");
ALTER TABLE "world_maps" ADD FOREIGN KEY ("world_id") REFERENCES "worlds" ("id");
ALTER TABLE "world_maps" ADD FOREIGN KEY ("map_id") REFERENCES "maps" ("id");

INSERT INTO world_maps (world_id, map_id)
SELECT
    module_type_id as world_id,
    map_id as map_id
FROM
    view_entities
WHERE
    type='map' AND
    module_type = 'world'
;


CREATE TABLE "world_images" (
    "world_id" int PRIMARY KEY,
    "header_img_id" int,
    "thumbnail_img_id" int,
    "avatar_img_id" int
);
ALTER TABLE "world_images" ADD FOREIGN KEY ("world_id") REFERENCES "worlds" ("id");
ALTER TABLE "world_images" ADD FOREIGN KEY ("header_img_id") REFERENCES "images" ("id");
ALTER TABLE "world_images" ADD FOREIGN KEY ("thumbnail_img_id") REFERENCES "images" ("id");
ALTER TABLE "world_images" ADD FOREIGN KEY ("avatar_img_id") REFERENCES "images" ("id");


INSERT INTO world_images (world_id, header_img_id, thumbnail_img_id, avatar_img_id)
SELECT
    world_id,
    header_img_id,
    thumbnail_img_id,
    avatar_img_id
FROM
    modules
WHERE
    module_type = 'world'
;


CREATE TABLE "world_menu" (
    "world_id" int NOT NULL,
    "menu_id" int NOT NULL
);
CREATE UNIQUE INDEX ON "world_menu" ("world_id", "menu_id");
ALTER TABLE "world_menu" ADD FOREIGN KEY ("world_id") REFERENCES "worlds" ("id");
ALTER TABLE "world_menu" ADD FOREIGN KEY ("menu_id") REFERENCES "menus" ("id");

INSERT INTO world_menu (world_id, menu_id)
SELECT
    world_id,
    menu_id
FROM
    modules
WHERE
    module_type = 'world'
;


-- ============================= RECREATE VIEWS

DROP VIEW view_connections_posts;
DROP FUNCTION get_worlds(boolean, integer[], varchar, varchar, int, int);
DROP VIEW view_worlds;
DROP VIEW view_modules;


ALTER TABLE "modules"
    DROP COLUMN menu_id,
    DROP COLUMN header_img_id,
    DROP COLUMN thumbnail_img_id,
    DROP COLUMN avatar_img_id
;


CREATE VIEW view_connections_menus AS
SELECT
    m.id as menu_id,
    COALESCE(wm.world_id, 0) as world_id,
    0 as quest_id,
    0 as character_id,
    0 as system_id
FROM
    menus m
    LEFT JOIN world_menu wm ON m.id = wm.menu_id
;

CREATE VIEW view_worlds AS
SELECT
    w.*,
    i_header.url as image_header,
    i_thumbnail.url as image_thumbnail,
    i_avatar.url as image_avatar,
    tags.tags AS tags,
    COALESCE(activity.activity_post_count, 0) AS activity_post_count,
    COALESCE(activity.activity_quest_count, 0) AS activity_quest_count,
    COALESCE(activity.activity_resource_count, 0) AS activity_resource_count,
    wm.menu_id as world_menu_id
FROM
    worlds w
    JOIN world_images wi ON w.id = wi.world_id
    JOIN world_menu wm ON w.id = wm.world_id
    LEFT JOIN (
        SELECT
            wa.world_id,
            cast(sum(wa.post_count) as integer) AS activity_post_count,
            cast(sum(wa.quest_count) as integer) AS activity_quest_count,
            cast(sum(wa.resource_count) as integer) AS activity_resource_count
        FROM
            world_activity wa
        WHERE
                wa.date >= (now() - interval '30 days')
        GROUP BY wa.world_id
    ) activity ON activity.world_id = w.id
    LEFT JOIN (
        SELECT
            m.world_id,
            cast(array_agg(t.id) as integer[]) AS tags
        FROM
            modules m
                JOIN module_tags mt ON m.id = mt.module_id
                LEFT JOIN module_type_tags_available t ON t.id = mt.tag_id
        WHERE world_id IS NOT NULL
        GROUP BY m.world_id
    ) tags ON tags.world_id = w.id
    LEFT JOIN images i_header on wi.header_img_id = i_header.id
    LEFT JOIN images i_thumbnail on wi.thumbnail_img_id = i_thumbnail.id
    LEFT JOIN images i_avatar on wi.avatar_img_id = i_avatar.id
;

CREATE VIEW view_connections_world_posts AS
SELECT
    wl.world_id as world_id,
    l.post_id as post_id,
    'locations' as helper_name,
    l.id as helper_id
FROM
    world_locations wl
    JOIN locations l ON wl.location_id = l.id JOIN posts p ON l.post_id = p.id
WHERE l.post_id IS NOT NULL

UNION ALL

SELECT
    wm.world_id as world_id,
    mip.post_id,
    'menu_item_posts',
    mip.menu_id
FROM
    world_menu wm
    JOIN menu_item_posts mip ON mip.menu_id = wm.menu_id

UNION ALL

SELECT
    wm.world_id,
    mi.description_post_id,
    'menu_items',
    mi.id
FROM
    world_menu wm
    JOIN menu_items mi ON mi.menu_id = wm.menu_id
WHERE mi.description_post_id IS NOT NULL

UNION ALL

SELECT
    id,
    description_post_id,
    'worlds',
    id
FROM
    worlds
WHERE description_post_id IS NOT NULL

UNION ALL

SELECT
    wm.world_id,
    e.post_id,
    'entities',
    e.id
FROM
    world_menu wm
    JOIN menu_items mi ON mi.menu_id = wm.menu_id
    JOIN get_recursive_entities(mi.entity_group_id) re ON 1 = 1
    JOIN entities e ON e.id = re.content_entity_id
WHERE e.post_id IS NOT NULL
;



CREATE OR REPLACE FUNCTION get_worlds(_is_public boolean, _tags integer[], _order_by varchar, _order_direction varchar, _limit int, _offset int)
    RETURNS SETOF view_worlds AS
$func$
BEGIN
    IF _order_by IS NULL THEN
        _order_by := 'created_at';
    END IF;

    IF _order_direction IS NULL OR (_order_direction <> 'ASC' AND _order_direction <> 'DESC') THEN
        _order_direction := 'DESC';
    END IF;

    RETURN QUERY EXECUTE format('
        SELECT * FROM view_worlds
        WHERE
            ($1 IS NULL OR public = $1) AND
            (array_length($2, 1) IS NULL OR tags @> $2)
        ORDER BY %I ' || _order_direction || '
        LIMIT $3
        OFFSET $4', _order_by)
        USING _is_public, _tags, _limit, _offset;
END
$func$  LANGUAGE plpgsql;