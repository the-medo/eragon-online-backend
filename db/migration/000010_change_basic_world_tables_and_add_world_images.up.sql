
ALTER TABLE "worlds"
    ADD COLUMN "description" varchar NOT NULL DEFAULT '',
    DROP COLUMN "img_id";

ALTER TABLE "world_admins"
    ADD COLUMN "is_main" boolean NOT NULL DEFAULT false;

--"description" varchar NOT NULL DEFAULT ''

CREATE TABLE "world_images" (
    "world_id" int PRIMARY KEY,
    "image_header" int,
    "image_avatar" int
);

ALTER TABLE "world_images" ADD FOREIGN KEY ("world_id") REFERENCES "worlds" ("id");

ALTER TABLE "world_images" ADD FOREIGN KEY ("image_header") REFERENCES "images" ("id");

ALTER TABLE "world_images" ADD FOREIGN KEY ("image_avatar") REFERENCES "images" ("id");

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

CREATE VIEW view_users AS
SELECT
    u.*,
    i.url as image_avatar
FROM
    users AS u
    LEFT JOIN images i ON u.img_id = i.id
;