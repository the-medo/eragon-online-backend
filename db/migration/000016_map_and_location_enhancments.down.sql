DROP PROCEDURE IF EXISTS assign_post_by_menu_id(INT, INT);

DROP INDEX world_map_pin_type_groups_world_id_map_pin_type_group_id_idx;
DROP INDEX world_maps_world_id_map_id_idx;
DROP INDEX world_locations_world_id_location_id_idx;

DROP TABLE IF EXISTS world_posts;

DROP VIEW view_connections_menus;
DROP VIEW view_connections_world_posts;
DROP VIEW view_locations;
CREATE VIEW view_locations AS
SELECT
    l.*,
    i.url as thumbnail_image_url
FROM
    locations l
        LEFT JOIN images i ON l.thumbnail_image_id = i.id
;

DROP PROCEDURE IF EXISTS delete_map(INT);
DROP PROCEDURE IF EXISTS delete_location(INT);