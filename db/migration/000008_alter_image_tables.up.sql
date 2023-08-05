CREATE TYPE "image_variant" AS ENUM (
    '100x100',
    '1200x800',
    '150x150',
    '1920x200',
    '200x200',
    '300x300',
    '30x30',
    'public'
);

-- Alter image_types table
ALTER TABLE "image_types"
    DROP COLUMN "width",
    DROP COLUMN "height",
    ADD COLUMN "variant" image_variant NOT NULL DEFAULT 'public',
    ALTER COLUMN "id" DROP IDENTITY;

DO $$
BEGIN
IF NOT EXISTS (
        SELECT constraint_name
        FROM information_schema.table_constraints
        WHERE table_name='image_types' AND constraint_name='image_types_id_uniq'
    ) THEN
ALTER TABLE "image_types"
    ADD CONSTRAINT image_types_id_uniq UNIQUE ("id");
END IF;
END
$$;

-- Alter images table
ALTER TABLE "images"
    DROP COLUMN "width",
    DROP COLUMN "height",
    ADD COLUMN "base_url" varchar NOT NULL DEFAULT '',
    ADD COLUMN "img_guid" uuid UNIQUE;

ALTER TABLE "images" RENAME COLUMN "address" TO "url";

COMMENT ON COLUMN "image_types"."variant" IS 'Variant name from cloudflare. ';

UPDATE "image_types" SET id = 200, variant = '200x200' WHERE name = 'User avatar';
UPDATE "image_types" SET id = 300, variant = '1920x300' WHERE name = 'World Header';
UPDATE "image_types" SET id = 400, variant = '200x200' WHERE name = 'World Avatar';
UPDATE "image_types" SET id = 500, variant = 'public' WHERE name = 'Location Image';
UPDATE "image_types" SET id = 600, variant = '300x300' WHERE name = 'Race Image';
UPDATE "image_types" SET id = 700, variant = '100x100' WHERE name = 'Item Image';
UPDATE "image_types" SET id = 800, variant = '100x100' WHERE name = 'Skill Image';
UPDATE "image_types" SET id = 900, variant = '150x150' WHERE name = 'Character Portrait';
UPDATE "image_types" SET id = 1000, variant = '1200x800' WHERE name = 'Map Image';
UPDATE "image_types" SET id = 1100, variant = 'public' WHERE name = 'Background Image';

INSERT INTO "image_types" ("id", "name", "variant", "description") VALUES (100, 'Default', 'public', 'Default image type.');