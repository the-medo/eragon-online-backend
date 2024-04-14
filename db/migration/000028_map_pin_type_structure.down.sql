ALTER TABLE "map_layers" ADD COLUMN "is_main" boolean NOT NULL DEFAULT false;
ALTER TABLE "map_pin_types" DROP COLUMN "is_default";
ALTER TABLE "map_pin_types" ADD COLUMN "section" varchar NOT NULL DEFAULT 'base';
