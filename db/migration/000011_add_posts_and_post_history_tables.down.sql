DROP VIEW IF EXISTS "view_posts";

DROP VIEW IF EXISTS "view_users";

ALTER TABLE "users" DROP COLUMN "introduction_post_id";

CREATE VIEW view_users AS
SELECT
    u.*,
    i.url as image_avatar
FROM
    users AS u
        LEFT JOIN images i ON u.img_id = i.id
;

DROP TABLE "post_history";
DROP TABLE "posts";
DROP TABLE "post_types";


-- ========= Overcomplicated removal of "original" image_variant ============
-- Step 1: Change "original" to "public" in "image_types" table
UPDATE image_types SET variant = 'public' WHERE variant = 'original';

-- Step 2: Create a new type, without "original"
CREATE TYPE "image_variant_new" AS ENUM (
    '100x100',
    '1200x800',
    '150x150',
    '1920x300',
    '200x200',
    '300x300',
    '30x30',
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