DROP PROCEDURE IF EXISTS delete_menu_item(INT);
DROP PROCEDURE IF EXISTS delete_entity_group(INT);
DROP PROCEDURE IF EXISTS delete_entity_group_content(INT, INT, INT);
DROP PROCEDURE IF EXISTS delete_entity_group_content_and_move(INT);

DROP FUNCTION IF EXISTS get_menu_id_of_entity_group(INT);
DROP PROCEDURE IF EXISTS move_entity_group_content(INT, INT);

DROP FUNCTION get_recursive_entities(_main_entity_group_id integer);


ALTER TABLE "menu_items" DROP COLUMN "entity_group_id";

DROP TABLE IF EXISTS "menu_item_entity_groups";
DROP TABLE IF EXISTS "entity_group_content";
DROP TABLE IF EXISTS "entity_groups";
DROP TABLE IF EXISTS "entities";

DROP TYPE IF EXISTS "entity_type";
