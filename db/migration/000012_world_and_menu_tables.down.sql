
DROP VIEW IF EXISTS view_worlds;

-- Drop `menu_item_posts` table
DROP TABLE menu_item_posts;

-- Drop `menu_items` table
DROP TABLE menu_items;

-- Drop `world_menu` table
DROP TABLE world_menu;

-- Drop `menus` table
DROP TABLE menus;

-- Drop `world_activity` table
DROP TABLE world_activity;

-- Drop `world_tags` table
DROP TABLE world_tags;

-- Drop `world_tags_available` table
DROP TABLE world_tags_available;

-- Revert changes in `world_images` table
ALTER TABLE world_images DROP CONSTRAINT world_images_avatar_img_id_fkey;
ALTER TABLE world_images DROP CONSTRAINT world_images_thumbnail_img_id_fkey;
ALTER TABLE world_images DROP COLUMN avatar_img_id;
ALTER TABLE world_images DROP COLUMN thumbnail_img_id;
ALTER TABLE world_images RENAME COLUMN header_img_id TO image_header;

-- Revert changes in `world_admins` table
COMMENT ON COLUMN "world_admins"."approved" IS '';
ALTER TABLE world_admins DROP COLUMN motivational_letter;
ALTER TABLE world_admins DROP COLUMN approved;
ALTER TABLE world_admins RENAME COLUMN super_admin TO is_main;

-- Revert changes in `worlds` table
ALTER TABLE worlds DROP COLUMN based_on;

CREATE TABLE "world_stats" (
   "world_id" int PRIMARY KEY,
   "final_content_rating" int NOT NULL DEFAULT 0,
   "final_activity" int NOT NULL DEFAULT 0
);

CREATE TABLE "world_stats_history" (
   "world_id" int,
   "final_content_rating" int NOT NULL DEFAULT 0,
   "final_activity" int NOT NULL DEFAULT 0,
   "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "world_stats" USING BTREE ("final_content_rating");
CREATE INDEX ON "world_stats" USING BTREE ("final_activity");
CREATE INDEX ON "world_stats_history" USING BTREE ("created_at");
ALTER TABLE "world_stats" ADD FOREIGN KEY ("world_id") REFERENCES "worlds" ("id");
ALTER TABLE "world_stats_history" ADD FOREIGN KEY ("world_id") REFERENCES "worlds" ("id");

CREATE VIEW view_worlds AS
SELECT
    w.*,
    i_avatar.url as image_avatar,
    i_header.url as image_header,
    ws.final_content_rating as rating,
    ws.final_activity as activity
FROM
    worlds w
    JOIN world_images wi ON w.id = wi.world_id
    JOIN world_stats ws ON w.id = wi.world_id
    LEFT JOIN images i_avatar on wi.image_avatar = i_avatar.id
    LEFT JOIN images i_header on wi.image_header = i_header.id
;