

DROP VIEW IF EXISTS "view_users";

ALTER TABLE "users" DROP COLUMN "introduction_post_id";

CREATE VIEW view_users AS
SELECT
    u.*,
    i.url as image_avatar
FROM
    users AS u
        LEFT JOIN images i ON u.img_id = i.id
;

DROP TABLE "post_history";
DROP TABLE "posts";
DROP TABLE "post_types";