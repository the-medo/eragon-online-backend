ALTER TABLE "images" DROP COLUMN "height";
ALTER TABLE "images" DROP COLUMN "width";
ALTER TABLE "menu_items" DROP COLUMN "entity_group_id";

DROP TABLE IF EXISTS "menu_item_entity_groups";
DROP TABLE IF EXISTS "entity_group_content";
DROP TABLE IF EXISTS "entity_groups";
DROP TABLE IF EXISTS "entities";

DROP TYPE IF EXISTS "entity_type";
