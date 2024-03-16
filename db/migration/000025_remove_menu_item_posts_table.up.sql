
DROP TABLE IF EXISTS menu_item_posts;

CREATE TYPE "delete_entity_group_content_action" AS ENUM (
    'delete_children',
    'move_children'
);