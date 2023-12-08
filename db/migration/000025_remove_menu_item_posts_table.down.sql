
CREATE TABLE "menu_item_posts" (
    "menu_id" int NOT NULL,
    "menu_item_id" int,
    "post_id" int NOT NULL,
    "position" int NOT NULL
);

CREATE UNIQUE INDEX ON "menu_item_posts" ("menu_item_id", "post_id");

ALTER TABLE "menu_item_posts" ADD FOREIGN KEY ("menu_id") REFERENCES "menus" ("id");

ALTER TABLE "menu_item_posts" ADD FOREIGN KEY ("menu_item_id") REFERENCES "menu_items" ("id");

ALTER TABLE "menu_item_posts" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id");
