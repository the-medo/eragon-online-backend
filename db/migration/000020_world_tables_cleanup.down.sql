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




ALTER TABLE "modules"
    DROP COLUMN menu_id,
    DROP COLUMN header_img_id,
    DROP COLUMN thumbnail_img_id,
    DROP COLUMN avatar_img_id
;
