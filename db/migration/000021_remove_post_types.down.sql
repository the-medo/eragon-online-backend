
DROP VIEW view_menu_item_posts;
DROP VIEW view_posts;


CREATE TABLE "post_types" (
    "id" int PRIMARY KEY,
    "name" varchar NOT NULL,
    "draftable" bool NOT NULL DEFAULT true,
    "privatable" bool NOT NULL DEFAULT false
);


INSERT INTO "post_types" ("id", "name", "draftable", "privatable")
VALUES
    (100, 'Universal', true, true),
    (200, 'Quest post', true, false),
    (300, 'World description', true, false),
    (400, 'Rule set description', true, false),
    (500, 'Quest description', true, false),
    (600, 'Character description', true, false),
    (700, 'News', true, false),
    (800, 'User introduction', true, false)
;

ALTER TABLE "posts" ADD COLUMN post_type_id integer NOT NULL default 100;
ALTER TABLE "post_history" ADD COLUMN post_type_id integer NOT NULL default 100;


ALTER TABLE "posts" ADD FOREIGN KEY ("post_type_id") REFERENCES "post_types" ("id");
ALTER TABLE "post_history" ADD FOREIGN KEY ("post_type_id") REFERENCES "post_types" ("id");


CREATE VIEW view_posts AS
SELECT
    p.*,
    pt.name as post_type_name,
    pt.draftable as post_type_draftable,
    pt.privatable as post_type_privatable,
    i.url as thumbnail_img_url,
    e.id as entity_id,
    e.module_id as module_id,
    e.module_type as module_type,
    e.module_type_id as module_type_id,
    e.tags as tags
FROM
    posts p
    JOIN post_types pt ON p.post_type_id = pt.id
    LEFT JOIN images i ON p.thumbnail_img_id = i.id
    LEFT JOIN view_entities e ON e.post_id = p.id
;

CREATE VIEW view_menu_item_posts AS
SELECT * FROM menu_item_posts mip JOIN view_posts vp ON mip.post_id = vp.id;