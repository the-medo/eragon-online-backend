ALTER TABLE "map_layers" ADD COLUMN position integer;

UPDATE map_layers SET position = 1 WHERE sublayer = true;

WITH layer_ranks AS (
    SELECT
        id,
        ROW_NUMBER() OVER(PARTITION BY map_id ORDER BY id) AS rn
    FROM map_layers WHERE sublayer = false
)
UPDATE map_layers
SET position = layer_ranks.rn + 1
FROM layer_ranks
WHERE map_layers.id = layer_ranks.id;

ALTER TABLE "map_layers" ALTER COLUMN position SET NOT NULL;

ALTER TABLE "map_layers" DROP COLUMN sublayer;
