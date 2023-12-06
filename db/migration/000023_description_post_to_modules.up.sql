
ALTER TABLE modules ADD COLUMN "description_post_id" int;
ALTER TABLE modules ADD CONSTRAINT modules_description_post_id_fkey FOREIGN KEY(description_post_id) REFERENCES posts(id);

UPDATE modules
SET description_post_id = worlds.description_post_id
FROM worlds
WHERE modules.world_id = worlds.id;

-- assuming that all worlds have "description_post_id" - which can be wrong, because it is nullable
ALTER TABLE modules ALTER COLUMN description_post_id SET NOT NULL;
ALTER TABLE worlds DROP COLUMN description_post_id;