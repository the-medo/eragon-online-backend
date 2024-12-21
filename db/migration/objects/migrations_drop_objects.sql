DROP PROCEDURE IF EXISTS update_map_layer_is_main(map_layer_id integer);
DROP PROCEDURE IF EXISTS move_menu_item(p_id integer, p_target_position integer);
DROP PROCEDURE IF EXISTS move_group_up(p_id integer);
DROP PROCEDURE IF EXISTS move_entity_group_content(p_id integer, p_target_position integer);
DROP VIEW IF EXISTS view_connections_posts;
DROP FUNCTION IF EXISTS get_recursive_entities(_main_entity_group_id integer);
DROP TYPE IF EXISTS get_recursive_entities_row;
DROP FUNCTION IF EXISTS get_menu_id_of_entity_group(_entity_group_id integer);
DROP PROCEDURE IF EXISTS delete_menu_item(_menu_item_id integer);
DROP PROCEDURE IF EXISTS delete_entity_group_content_and_move(_id integer);
DROP PROCEDURE IF EXISTS delete_entity_group_content(_id integer, _delete_type delete_entity_group_content_action);
DROP PROCEDURE IF EXISTS delete_entity_group_content(_id integer, _entity_id integer, _entity_group_id integer, _delete_type delete_entity_group_content_action);
DROP PROCEDURE IF EXISTS delete_entity_group_content(_id integer, _entity_id integer, _entity_group_id integer);
DROP PROCEDURE IF EXISTS delete_entity_group(_entity_group_id integer, _delete_type delete_entity_group_content_action);
DROP PROCEDURE IF EXISTS delete_entity_group(_entity_group_id integer);
DROP PROCEDURE IF EXISTS delete_map(p_map_id integer);
DROP PROCEDURE IF EXISTS delete_location(p_location_id integer);
DROP FUNCTION IF EXISTS get_worlds(boolean, integer[], varchar, varchar, int, int);
DROP FUNCTION IF EXISTS get_systems(boolean, integer[], varchar, varchar, int, int);
DROP FUNCTION IF EXISTS get_quests CASCADE;
DROP FUNCTION IF EXISTS get_characters(boolean, integer[], int, int, varchar, varchar, int, int);
DROP FUNCTION IF EXISTS get_posts(_is_private boolean, _is_draft boolean, _tags integer[], _user_id integer, _module_id integer, _module_type module_type, _order_by varchar, _order_direction varchar, _limit int, _offset int);
DROP FUNCTION IF EXISTS get_images(_tags integer[], _width integer, _height integer, _user_id integer, _module_id integer, _module_type module_type, _order_by varchar, _order_direction varchar, _limit int, _offset int);
DROP FUNCTION IF EXISTS get_locations(_tags integer[], _module_id integer, _module_type module_type, _order_by varchar, _order_direction varchar, _limit int, _offset int);
DROP FUNCTION IF EXISTS get_maps(_tags integer[], _module_id integer, _module_type module_type, _order_by varchar, _order_direction varchar, _limit int, _offset int);

DROP VIEW IF EXISTS view_worlds;            -- depends on view_modules
DROP VIEW IF EXISTS view_systems;           -- depends on view_modules
DROP VIEW IF EXISTS view_quests;            -- depends on view_modules
DROP VIEW IF EXISTS view_characters;        -- depends on view_modules
DROP VIEW IF EXISTS view_modules;
DROP VIEW IF EXISTS view_locations;         -- depends on view_entities
DROP VIEW IF EXISTS view_images;            -- depends on view_entities
DROP VIEW IF EXISTS view_maps;              -- depends on view_entities
DROP VIEW IF EXISTS view_posts;             -- depends on view_entities
DROP VIEW IF EXISTS view_entities;
DROP VIEW IF EXISTS view_map_layers;
DROP VIEW IF EXISTS view_map_pins;
DROP VIEW IF EXISTS view_menus;
DROP VIEW IF EXISTS view_module_admins;
DROP VIEW IF EXISTS view_module_type_tags_available;
DROP VIEW IF EXISTS view_users;


