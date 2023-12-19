
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