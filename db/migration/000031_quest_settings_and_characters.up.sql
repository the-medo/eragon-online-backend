CREATE TYPE "quest_status" AS ENUM (
  'unknown',
  'not_started',
  'in_progress',
  'finished_completed',
  'finished_not_completed'
);

CREATE TABLE "quest_characters" (
    "quest_id" int NOT NULL,
    "character_id" int NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "approved" int NOT NULL,
    "motivational_letter" varchar NOT NULL
);

CREATE TABLE "quest_settings" (
    "quest_id" int NOT NULL,
    "status" quest_status NOT NULL DEFAULT 'not_started',
    "can_join" boolean NOT NULL DEFAULT false
);

INSERT INTO quest_settings (quest_id) SELECT id FROM quests;

CREATE UNIQUE INDEX ON "quest_characters" ("quest_id", "character_id");

CREATE UNIQUE INDEX ON "quest_settings" ("quest_id");

COMMENT ON COLUMN "quest_characters"."approved" IS '0 = NO, 1 = YES, 2 = PENDING';

ALTER TABLE "quest_characters" ADD FOREIGN KEY ("quest_id") REFERENCES "quests" ("id");

ALTER TABLE "quest_characters" ADD FOREIGN KEY ("character_id") REFERENCES "characters" ("id");

ALTER TABLE "quest_settings" ADD FOREIGN KEY ("quest_id") REFERENCES "quests" ("id");
