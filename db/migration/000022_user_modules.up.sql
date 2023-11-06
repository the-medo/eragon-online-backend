CREATE TABLE "user_modules" (
    "user_id" int NOT NULL,
    "module_id" int NOT NULL,
    "admin" bool NOT NULL,
    "favorite" bool NOT NULL,
    "following" bool NOT NULL,
    "entity_notifications" entity_type[]
);

CREATE UNIQUE INDEX ON "user_modules" ("user_id", "module_id");

ALTER TABLE "user_modules" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_modules" ADD FOREIGN KEY ("module_id") REFERENCES "modules" ("id");

INSERT INTO user_modules (user_id, module_id, admin, favorite, following, entity_notifications)
SELECT user_id, module_id, true, true, true, '{}'::entity_type[] FROM module_admins;