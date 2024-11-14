-- Add new columns to "quests" table
ALTER TABLE "quests"
    ADD COLUMN "name" varchar UNIQUE NOT NULL,
ADD COLUMN "public" boolean NOT NULL DEFAULT false,
ADD COLUMN "created_at" timestamptz NOT NULL DEFAULT now(),
ADD COLUMN "short_description" varchar NOT NULL DEFAULT '',
ADD COLUMN "world_id" int,
ADD COLUMN "system_id" int;

-- Add new columns to "systems" table
ALTER TABLE "systems"
    ADD COLUMN "name" varchar UNIQUE NOT NULL,
ADD COLUMN "based_on" varchar NOT NULL,
ADD COLUMN "public" boolean NOT NULL DEFAULT false,
ADD COLUMN "created_at" timestamptz NOT NULL DEFAULT now(),
ADD COLUMN "short_description" varchar NOT NULL DEFAULT '';

-- Add new columns to "characters" table
ALTER TABLE "characters"
    ADD COLUMN "name" varchar UNIQUE NOT NULL,
ADD COLUMN "public" boolean NOT NULL DEFAULT false,
ADD COLUMN "created_at" timestamptz NOT NULL DEFAULT now(),
ADD COLUMN "short_description" varchar NOT NULL DEFAULT '',
ADD COLUMN "world_id" int,
ADD COLUMN "system_id" int;

-- Add foreign key constraint for "world_id" and "system_id" in "quests" table
ALTER TABLE "quests"
    ADD CONSTRAINT "fk_quests_world_id" FOREIGN KEY ("world_id") REFERENCES "worlds" ("id"),
ADD CONSTRAINT "fk_quests_system_id" FOREIGN KEY ("system_id") REFERENCES "systems" ("id");

-- Add foreign key constraint for "world_id" and "system_id" in "characters" table
ALTER TABLE "characters"
    ADD CONSTRAINT "fk_characters_world_id" FOREIGN KEY ("world_id") REFERENCES "worlds" ("id"),
ADD CONSTRAINT "fk_characters_system_id" FOREIGN KEY ("system_id") REFERENCES "systems" ("id");