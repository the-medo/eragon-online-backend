DROP VIEW IF EXISTS view_users;
DROP VIEW IF EXISTS view_worlds;

DROP TABLE "world_stats_history";
DROP TABLE "world_stats";
DROP TABLE "world_images";

ALTER TABLE "world_admins" DROP COLUMN "is_main";

ALTER TABLE "worlds"
    ADD COLUMN "img_id" int,
    DROP COLUMN "description";

ALTER TABLE "worlds" ADD FOREIGN KEY ("img_id") REFERENCES "images" ("id");
