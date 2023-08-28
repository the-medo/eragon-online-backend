DROP VIEW view_posts;
CREATE VIEW view_posts AS
SELECT
    p.*,
    pt.name as post_type_name,
    pt.draftable as post_type_draftable,
    pt.privatable as post_type_privatable
FROM
    posts p
    JOIN post_types pt ON p.post_type_id = pt.id
;

-- Alter the `post_history` table
ALTER TABLE "post_history"
    DROP COLUMN "description",
    DROP COLUMN "thumbnail_img_id";

-- Alter the `posts` table
ALTER TABLE "posts"
    DROP COLUMN "description",
    DROP COLUMN "thumbnail_img_id";

DROP VIEW IF EXISTS view_menus;
DROP PROCEDURE IF EXISTS move_menu_item(INT, INT);
DROP PROCEDURE IF EXISTS move_group_up(INT);

DELETE FROM "image_types" WHERE id IN (1200, 1300);

-- ========= Overcomplicated removal of "original" image_variant ============
-- Step 1: Change "original" to "public" in "image_types" table
-- UPDATE image_types SET variant = '1200x800' WHERE variant IN ('300x200', '600x400');
-- UPDATE image_types SET variant = '300x300' WHERE variant = '400x600';
-- UPDATE image_types SET variant = '200x200' WHERE variant = '200x300';

-- Step 2: Create a new type, without '300x200', '600x400', '400x600', '200x300'
CREATE TYPE "image_variant_new" AS ENUM (
    '100x100',
    '1200x800',
    '150x150',
    '1920x300',
    '200x200',
    '300x300',
    '30x30',
    'original',
    'public'
);


-- Step 3: Update the table to use the new type

-- Step 3.1: Remove the default value from the column
ALTER TABLE image_types ALTER COLUMN variant DROP DEFAULT;

-- Step 3.2: Change the column type to the new type
BEGIN;
ALTER TABLE image_types ALTER COLUMN variant TYPE image_variant_new USING variant::text::image_variant_new;
COMMIT;

-- Step 3.3: Add the default value back (with the new type)
ALTER TABLE image_types ALTER COLUMN variant SET DEFAULT 'public';


-- Step 4: Delete the old type
DROP TYPE image_variant;

-- Step 5: Rename the new type to the old type's name
ALTER TYPE image_variant_new RENAME TO image_variant;

-- =======================================================================


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
ALTER TABLE worlds DROP COLUMN description_post_id;
ALTER TABLE worlds RENAME COLUMN short_description TO description;
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