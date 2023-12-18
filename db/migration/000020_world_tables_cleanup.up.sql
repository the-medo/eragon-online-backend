
CREATE TABLE "module_map_pin_type_groups" (
    "module_id" int NOT NULL,
    "map_pin_type_group_id" int NOT NULL
);

CREATE UNIQUE INDEX ON "module_map_pin_type_groups" ("module_id", "map_pin_type_group_id");
ALTER TABLE "module_map_pin_type_groups" ADD FOREIGN KEY ("module_id") REFERENCES "modules" ("id");
ALTER TABLE "module_map_pin_type_groups" ADD FOREIGN KEY ("map_pin_type_group_id") REFERENCES "map_pin_type_group" ("id");

INSERT INTO module_map_pin_type_groups (module_id, map_pin_type_group_id)
SELECT
    m.id as module_id,
    wmpg.map_pin_type_group_id as map_pin_type_group_id
FROM
    modules m
    JOIN world_map_pin_type_groups wmpg ON wmpg.world_id = m.world_id;



ALTER TABLE "modules"
    ADD COLUMN menu_id integer,
    ADD COLUMN header_img_id integer,
    ADD COLUMN thumbnail_img_id integer,
    ADD COLUMN avatar_img_id integer
;

ALTER TABLE "modules" ADD FOREIGN KEY ("menu_id") REFERENCES "menus" ("id");
ALTER TABLE "modules" ADD FOREIGN KEY ("header_img_id") REFERENCES "images" ("id");
ALTER TABLE "modules" ADD FOREIGN KEY ("thumbnail_img_id") REFERENCES "images" ("id");
ALTER TABLE "modules" ADD FOREIGN KEY ("avatar_img_id") REFERENCES "images" ("id");

WITH module_menus AS (
    SELECT
        m.id as module_id,
        wm.menu_id as menu_id,
        wi.header_img_id as header_img_id,
        wi.thumbnail_img_id as thumbnail_img_id,
        wi.avatar_img_id as avatar_img_id
    FROM
        modules m
        JOIN worlds w ON w.id = m.world_id
        JOIN world_menu wm ON wm.world_id = w.id
        JOIN world_images wi ON wi.world_id = w.id
)
UPDATE modules
    SET
        menu_id = module_menus.menu_id,
        header_img_id = module_menus.header_img_id,
        thumbnail_img_id = module_menus.thumbnail_img_id,
        avatar_img_id = module_menus.avatar_img_id
    FROM module_menus
    WHERE module_menus.module_id = modules.id
;

DROP VIEW view_connections_world_posts;
DROP VIEW view_connections_menus;
DROP FUNCTION get_worlds(boolean, integer[], varchar, varchar, int, int);
DROP VIEW view_worlds;

DROP TABLE world_menu;
DROP TABLE world_images;
DROP TABLE world_maps;
DROP TABLE world_locations;
DROP TABLE world_posts;
DROP TABLE world_map_pin_type_groups;
DROP TABLE world_activity;


CREATE TABLE "module_admins" (
    "module_id" int NOT NULL,
    "user_id" int NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "super_admin" boolean NOT NULL DEFAULT false,
    "approved" int NOT NULL,
    "motivational_letter" varchar NOT NULL,
    "allowed_entity_types" entity_type[],
    "allowed_menu" boolean NOT NULL DEFAULT false
);

CREATE UNIQUE INDEX ON "module_admins" ("module_id", "user_id");
COMMENT ON COLUMN "module_admins"."approved" IS '0 = NO, 1 = YES, 2 = PENDING';
ALTER TABLE "module_admins" ADD FOREIGN KEY ("module_id") REFERENCES "modules" ("id");
ALTER TABLE "module_admins" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

INSERT INTO module_admins (module_id, user_id, super_admin, approved, motivational_letter, allowed_entity_types, allowed_menu)
SELECT
    m.id as module_id,
    wa.user_id as user_id,
    wa.super_admin as super_admin,
    wa.approved as approved,
    wa.motivational_letter as motivational_letter,
    ARRAY[]::entity_type[] as allowed_entity_types,
    wa.super_admin as allowed_menu
FROM
    modules m
    JOIN world_admins wa ON wa.world_id = m.world_id
WHERE
    m.module_type = 'world'
;

DROP TABLE world_admins;
