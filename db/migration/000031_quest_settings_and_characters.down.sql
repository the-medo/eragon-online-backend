ALTER TABLE "quests"
    DROP COLUMN IF EXISTS "status",
    DROP COLUMN IF EXISTS "can_join";

DROP TABLE IF EXISTS "quest_characters";
DROP TYPE IF EXISTS "quest_status";
