DELETE FROM "image_types";
DELETE FROM "roles";

ALTER TABLE "user_roles" DROP COLUMN "created_at";
ALTER TABLE "images" DROP COLUMN "created_at";
ALTER TABLE "image_types" DROP COLUMN "description";
ALTER TABLE "roles" DROP COLUMN "description";