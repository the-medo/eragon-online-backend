ALTER TABLE "maps" ADD COLUMN "user_id" int NOT NULL DEFAULT 0;
ALTER TABLE "maps" ADD COLUMN "created_at" timestamptz NOT NULL DEFAULT (now());
ALTER TABLE "maps" ADD COLUMN "last_updated_at" timestamptz;
ALTER TABLE "maps" ADD COLUMN "last_updated_user_id" int;
ALTER TABLE "maps" ADD COLUMN "is_private" boolean NOT NULL DEFAULT false;
ALTER TABLE "maps" RENAME COLUMN "name" TO "title";

UPDATE maps
SET user_id = ma.user_id
FROM
    entities e
        JOIN module_admins ma ON ma.module_id = e.module_id
WHERE
    ma.super_admin = true AND
    maps.id = e.map_id;

ALTER TABLE "maps" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "maps" ADD FOREIGN KEY ("last_updated_user_id") REFERENCES "users" ("id");
