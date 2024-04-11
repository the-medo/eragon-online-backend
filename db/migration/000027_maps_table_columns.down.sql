ALTER TABLE "maps" DROP COLUMN "user_id";
ALTER TABLE "maps" DROP COLUMN "created_at";
ALTER TABLE "maps" DROP COLUMN "last_updated_at";
ALTER TABLE "maps" DROP COLUMN "last_updated_user_id";
ALTER TABLE "maps" DROP COLUMN "is_private";
ALTER TABLE "maps" RENAME COLUMN "title" TO "name";