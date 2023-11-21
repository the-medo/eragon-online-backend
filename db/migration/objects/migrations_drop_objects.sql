DROP PROCEDURE update_map_layer_is_main(map_layer_id integer);
DROP PROCEDURE move_menu_item_post(mi_id integer, p_id integer, p_target_position integer);
DROP PROCEDURE move_menu_item(p_id integer, p_target_position integer);
DROP PROCEDURE move_group_up(p_id integer);
DROP PROCEDURE move_entity_group_content(p_id integer, p_target_position integer);
DROP VIEW view_connections_posts;
DROP FUNCTION get_recursive_entities(_main_entity_group_id integer);
DROP FUNCTION get_menu_id_of_entity_group(_entity_group_id integer);
DROP PROCEDURE delete_menu_item(_menu_item_id integer);
DROP PROCEDURE delete_entity_group_content_and_move(_id integer);
DROP PROCEDURE delete_entity_group_content(_id integer, _entity_id integer, _entity_group_id integer);
DROP PROCEDURE delete_entity_group(_entity_group_id integer);
DROP PROCEDURE delete_map(p_map_id integer);
DROP PROCEDURE delete_location(p_location_id integer);
DROP FUNCTION get_worlds(boolean, integer[], varchar, varchar, int, int);

DROP VIEW view_worlds;
DROP VIEW view_modules;
DROP VIEW view_locations; -- depends on view_entities
DROP VIEW view_images; -- depends on view_entities
DROP VIEW view_maps;-- depends on view_entities
DROP VIEW view_menu_item_posts; -- depends on view_posts
DROP VIEW view_posts; -- depends on view_entities
DROP VIEW view_entities;
DROP VIEW view_map_layers;
DROP VIEW view_map_pins;
DROP VIEW view_menus;
DROP VIEW view_module_admins;
DROP VIEW view_module_type_tags_available;
DROP VIEW view_users;


