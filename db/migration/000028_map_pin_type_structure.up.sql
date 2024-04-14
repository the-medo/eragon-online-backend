
ALTER TABLE "map_layers" DROP COLUMN "is_main";
ALTER TABLE "map_pin_types" DROP COLUMN "section";
ALTER TABLE "map_pin_types" ADD COLUMN "is_default" boolean NOT NULL DEFAULT false;

UPDATE map_pin_types
SET is_default = true
WHERE id IN (
    SELECT
        MIN(mpt.id) as map_pin_type_id
    FROM
        module_map_pin_type_groups mmpt
        JOIN map_pin_type_group mptg ON mptg.id = mmpt.map_pin_type_group_id
        JOIN map_pin_types mpt ON mpt.map_pin_type_group_id = mptg.id
    GROUP BY
        mmpt.module_id
);
