
ALTER TABLE worlds ADD COLUMN "description_post_id" int;
ALTER TABLE worlds ADD CONSTRAINT worlds_description_post_id_fkey FOREIGN KEY(description_post_id) REFERENCES posts(id);

UPDATE worlds
SET description_post_id = modules.description_post_id
FROM modules
WHERE worlds.id = modules.world_id;

ALTER TABLE modules DROP COLUMN description_post_id;