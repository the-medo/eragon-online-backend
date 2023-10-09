DROP PROCEDURE IF EXISTS delete_entity_group(INT);
DROP PROCEDURE IF EXISTS delete_entity_group_content(INT, INT, INT);
DROP PROCEDURE IF EXISTS delete_entity_group_content_and_move(INT);

DROP FUNCTION IF EXISTS get_menu_id_of_entity_group(INT);
DROP PROCEDURE IF EXISTS move_entity_group_content(INT, INT);
DROP PROCEDURE IF EXISTS move_menu_entity_groups(INT, INT, INT);