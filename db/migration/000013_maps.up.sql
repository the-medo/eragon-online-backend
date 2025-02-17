CREATE TYPE "pin_shape" AS ENUM (
    'none',
    'square',
    'triangle',
    'pin',
    'circle',
    'hexagon',
    'octagon',
    'star',
    'diamond',
    'pentagon',
    'heart',
    'cloud'
);

CREATE TABLE "maps" (
    "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    "name" varchar NOT NULL,
    "type" varchar,
    "description" varchar,
    "width" int NOT NULL,
    "height" int NOT NULL,
    "thumbnail_image_id" int
);

CREATE TABLE "locations" (
     "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
     "name" varchar NOT NULL,
     "description" varchar,
     "post_id" int,
     "thumbnail_image_id" int
);

CREATE TABLE "map_layers" (
    "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    "name" varchar NOT NULL,
    "map_id" int NOT NULL,
    "image_id" int NOT NULL,
    "is_main" bool NOT NULL DEFAULT false,
    "enabled" bool NOT NULL DEFAULT true,
    "sublayer" bool NOT NULL DEFAULT false
);

CREATE TABLE "map_pin_types" (
    "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    "map_id" int NOT NULL,
    "shape" pin_shape NOT NULL DEFAULT 'pin',
    "background_color" varchar,
    "border_color" varchar,
    "icon_color" varchar,
    "icon" varchar,
    "icon_size" int,
    "width" int
);

CREATE TABLE "map_pins" (
    "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    "name" varchar NOT NULL,
    "map_id" int NOT NULL,
    "map_pin_type_id" int,
    "location_id" int,
    "map_layer_id" int,
    "x" int NOT NULL,
    "y" int NOT NULL
);

CREATE TABLE "world_maps" (
    "world_id" int NOT NULL,
    "map_id" int NOT NULL
);

CREATE TABLE "world_locations" (
    "world_id" int NOT NULL,
    "location_id" int NOT NULL
);

ALTER TABLE "maps" ADD FOREIGN KEY ("thumbnail_image_id") REFERENCES "images" ("id");

ALTER TABLE "locations" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id");

ALTER TABLE "locations" ADD FOREIGN KEY ("thumbnail_image_id") REFERENCES "images" ("id");

ALTER TABLE "map_layers" ADD FOREIGN KEY ("map_id") REFERENCES "maps" ("id");

ALTER TABLE "map_layers" ADD FOREIGN KEY ("image_id") REFERENCES "images" ("id");

ALTER TABLE "map_pins" ADD FOREIGN KEY ("map_id") REFERENCES "maps" ("id");

ALTER TABLE "map_pins" ADD FOREIGN KEY ("location_id") REFERENCES "locations" ("id");

ALTER TABLE "map_pins" ADD FOREIGN KEY ("map_layer_id") REFERENCES "map_layers" ("id");

ALTER TABLE "map_pins" ADD FOREIGN KEY ("map_pin_type_id") REFERENCES "map_pin_types" ("id");

ALTER TABLE "map_pin_types" ADD FOREIGN KEY ("map_id") REFERENCES "maps" ("id");

ALTER TABLE "world_maps" ADD FOREIGN KEY ("world_id") REFERENCES "worlds" ("id");

ALTER TABLE "world_maps" ADD FOREIGN KEY ("map_id") REFERENCES "maps" ("id");

ALTER TABLE "world_locations" ADD FOREIGN KEY ("world_id") REFERENCES "worlds" ("id");

ALTER TABLE "world_locations" ADD FOREIGN KEY ("location_id") REFERENCES "locations" ("id");