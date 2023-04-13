ALTER TABLE "roles" ADD COLUMN "description" text NOT NULL;
ALTER TABLE "image_types" ADD COLUMN "description" text NOT NULL;
ALTER TABLE "images" ADD COLUMN "created_at" timestamptz NOT NULL DEFAULT (now());
ALTER TABLE "user_roles" ADD COLUMN "created_at" timestamptz NOT NULL DEFAULT (now());

INSERT INTO "roles" ("name", "description")
VALUES
    ('admin', 'Can do everything'),
    ('moderator', 'Can moderate the site')
;

INSERT INTO "image_types" ("name", "width", "height", "description")
VALUES
    ('User avatar', 150, 150, 'Image that represents player in chat and profile. Recommended size: 150x150px.'),
    ('World Header', 1000, 200, 'A banner image that represents the world''s theme, setting, or key locations, used as a thumbnail or cover photo. Recommended size: 1000x200px.'),
    ('World Avatar', 200, 200, 'A square image that serves as the world''s icon on the website, used to represent the world in various listings and previews. Recommended size: 200x200px.'),
    ('Location Image', 800, 600, 'Images of specific locations within the world, such as cities, dungeons, or landscapes, which can be displayed on the world''s page or used in the world-building tools. Recommended size: 800x600px.'),
    ('Race Image', 300, 300, 'Visual representations of the different playable races, which can be displayed during character creation or on the world''s page. Recommended size: 300x300px.'),
    ('Item Image', 100, 100, 'Images of unique or custom items available in the world, used for inventory management, item descriptions, or in-game shops. Recommended size: 100x100px.'),
    ('Skill Image', 100, 100, 'Images of unique or custom skill available in the world, used for character creation. Recommended size: 100x100px.'),
    ('Character Portrait', 150, 150, 'Images that represent player characters or key NPCs, which can be used in character sheets, chat interfaces, or as avatars during gameplay. Recommended size: 150x150px.'),
    ('Map Image', 1200, 800, 'Visuals of world maps, regional maps, or dungeon maps, which can be displayed on the world''s page or used as a reference during gameplay. Recommended size: 1200x800px or larger, depending on the level of detail required.'),
    ('Background Image', 1920, 1080, 'Atmospheric images or patterns that can be used as background visuals for the world''s page or in-game interfaces, helping to set the mood or theme. Recommended size: 1920x1080px or larger, depending on the desired resolution.')
;

