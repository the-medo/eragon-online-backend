ALTER TABLE "map_layers" ADD COLUMN sublayer boolean NOT NULL DEFAULT false;
UPDATE map_layers SET sublayer = true WHERE position = 1;

ALTER TABLE "map_layers" DROP COLUMN position;
