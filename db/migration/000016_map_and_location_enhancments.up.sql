
CREATE TABLE "world_posts" (
    "world_id" int NOT NULL,
    "post_id" int NOT NULL
);
ALTER TABLE "world_posts" ADD FOREIGN KEY ("world_id") REFERENCES "worlds" ("id");
ALTER TABLE "world_posts" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id");
CREATE UNIQUE INDEX ON "world_posts" ("world_id", "post_id");

-- INSERT INTO world_posts (world_id, post_id)
-- SELECT world_id, post_id FROM view_connections_world_posts
-- ON CONFLICT (world_id, post_id) DO NOTHING
-- ;

CREATE UNIQUE INDEX ON "world_map_pin_type_groups" ("world_id", "map_pin_type_group_id");
CREATE UNIQUE INDEX ON "world_maps" ("world_id", "map_id");
CREATE UNIQUE INDEX ON "world_locations" ("world_id", "location_id");
