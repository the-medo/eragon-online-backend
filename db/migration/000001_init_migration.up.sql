CREATE TABLE "users" (
  "id" integer PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "avatar" varchar,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "roles" (
  "id" integer PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL
);

CREATE TABLE "user_roles" (
  "user_id" integer NOT NULL,
  "role_id" integer NOT NULL
);

CREATE TABLE "chat" (
  "id" integer PRIMARY KEY,
  "user_id" integer NOT NULL,
  "text" varchar NOT NULL,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "worlds" (
  "id" integer PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "avatar" varchar,
  "public" boolean NOT NULL DEFAULT false
);

CREATE TABLE "world_admins" (
  "world_id" integer,
  "user_id" integer
);

CREATE TABLE "races" (
  "id" integer PRIMARY KEY,
  "world_id" integer NOT NULL,
  "name" varchar NOT NULL,
  "avatar" varchar,
  "is_playable" boolean NOT NULL DEFAULT true
);

CREATE TABLE "properties" (
  "id" integer PRIMARY KEY,
  "world_id" integer NOT NULL,
  "name" varchar NOT NULL
);

CREATE TABLE "race_properties" (
  "race_id" integer NOT NULL,
  "property_id" integer NOT NULL,
  "min_value" integer NOT NULL,
  "max_value" integer NOT NULL
);

CREATE TABLE "characters" (
  "id" integer PRIMARY KEY,
  "user_id" integer NOT NULL,
  "world_id" integer NOT NULL,
  "race_id" integer NOT NULL,
  "name" varchar NOT NULL,
  "level" integer NOT NULL,
  "experience" integer NOT NULL,
  "skill_points" integer NOT NULL
);

CREATE TABLE "character_properties" (
  "character_id" integer NOT NULL,
  "property_id" integer NOT NULL,
  "value" integer NOT NULL
);

CREATE TABLE "skills" (
  "id" integer PRIMARY KEY,
  "world_id" integer NOT NULL,
  "avatar" varchar,
  "name" varchar NOT NULL,
  "max_level" integer NOT NULL
);

CREATE TABLE "skill_requirements" (
  "id" integer PRIMARY KEY,
  "skill_id" integer NOT NULL,
  "level" integer,
  "race_id" integer
);

CREATE TABLE "skill_requirement_skills" (
  "skill_requirement_id" integer NOT NULL,
  "skill_id" integer NOT NULL,
  "level" integer NOT NULL
);

CREATE TABLE "skill_requirement_properties" (
  "skill_requirement_id" integer NOT NULL,
  "property_id" integer NOT NULL,
  "value" integer NOT NULL
);

CREATE UNIQUE INDEX ON "user_roles" ("role_id", "user_id");

CREATE UNIQUE INDEX ON "world_admins" ("world_id", "user_id");

CREATE UNIQUE INDEX ON "races" ("world_id", "name");

CREATE UNIQUE INDEX ON "properties" ("world_id", "name");

CREATE UNIQUE INDEX ON "race_properties" ("race_id", "property_id");

CREATE UNIQUE INDEX ON "characters" ("world_id", "name");

CREATE UNIQUE INDEX ON "character_properties" ("character_id", "property_id");

ALTER TABLE "user_roles" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_roles" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");

ALTER TABLE "chat" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "world_admins" ADD FOREIGN KEY ("world_id") REFERENCES "worlds" ("id");

ALTER TABLE "world_admins" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "races" ADD FOREIGN KEY ("world_id") REFERENCES "worlds" ("id");

ALTER TABLE "properties" ADD FOREIGN KEY ("world_id") REFERENCES "worlds" ("id");

ALTER TABLE "race_properties" ADD FOREIGN KEY ("race_id") REFERENCES "races" ("id");

ALTER TABLE "race_properties" ADD FOREIGN KEY ("property_id") REFERENCES "properties" ("id");

ALTER TABLE "characters" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "characters" ADD FOREIGN KEY ("race_id") REFERENCES "races" ("id");

ALTER TABLE "characters" ADD FOREIGN KEY ("world_id") REFERENCES "worlds" ("id");

ALTER TABLE "character_properties" ADD FOREIGN KEY ("character_id") REFERENCES "characters" ("id");

ALTER TABLE "character_properties" ADD FOREIGN KEY ("property_id") REFERENCES "properties" ("id");

ALTER TABLE "skills" ADD FOREIGN KEY ("world_id") REFERENCES "worlds" ("id");

ALTER TABLE "skill_requirements" ADD FOREIGN KEY ("skill_id") REFERENCES "skills" ("id");

ALTER TABLE "skill_requirements" ADD FOREIGN KEY ("race_id") REFERENCES "races" ("id");

ALTER TABLE "skill_requirement_skills" ADD FOREIGN KEY ("skill_requirement_id") REFERENCES "skill_requirements" ("id");

ALTER TABLE "skill_requirement_skills" ADD FOREIGN KEY ("skill_id") REFERENCES "skills" ("id");

ALTER TABLE "skill_requirement_properties" ADD FOREIGN KEY ("skill_requirement_id") REFERENCES "skill_requirements" ("id");

ALTER TABLE "skill_requirement_properties" ADD FOREIGN KEY ("property_id") REFERENCES "properties" ("id");
