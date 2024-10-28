-- Remove foreign key constraints from "quests" table
ALTER TABLE "quests" DROP CONSTRAINT IF EXISTS "fk_quests_world_id";
ALTER TABLE "quests" DROP CONSTRAINT IF EXISTS "fk_quests_system_id";

-- Remove foreign key constraints from "characters" table
ALTER TABLE "characters" DROP CONSTRAINT IF EXISTS "fk_characters_world_id";
ALTER TABLE "characters" DROP CONSTRAINT IF EXISTS "fk_characters_system_id";

-- Remove columns from "quests" table
ALTER TABLE "quests"
DROP COLUMN IF EXISTS "name",
DROP COLUMN IF EXISTS "public",
DROP COLUMN IF EXISTS "created_at",
DROP COLUMN IF EXISTS "short_description",
DROP COLUMN IF EXISTS "world_id",
DROP COLUMN IF EXISTS "system_id";

-- Remove columns from "systems" table
ALTER TABLE "systems"
DROP COLUMN IF EXISTS "name",
DROP COLUMN IF EXISTS "based_on",
DROP COLUMN IF EXISTS "public",
DROP COLUMN IF EXISTS "created_at",
DROP COLUMN IF EXISTS "short_description";

-- Remove columns from "characters" table
ALTER TABLE "characters"
DROP COLUMN IF EXISTS "name",
DROP COLUMN IF EXISTS "public",
DROP COLUMN IF EXISTS "created_at",
DROP COLUMN IF EXISTS "short_description",
DROP COLUMN IF EXISTS "world_id",
DROP COLUMN IF EXISTS "system_id";
