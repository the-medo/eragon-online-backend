

ALTER TABLE "images" ADD COLUMN "user_id" int NOT NULL DEFAULT 0;

DO
$$
    DECLARE
        admin_id INT;
    BEGIN
        SELECT
            u.id INTO admin_id
        FROM
            users AS u
            JOIN user_roles ur ON u.id = ur.user_id
            JOIN roles r ON ur.role_id = r.id
        WHERE
            r.name = 'admin'
        LIMIT 1;

        -- If no admin user id was found, get the id of the first user
        IF admin_id IS NULL THEN
            SELECT id INTO admin_id FROM users ORDER BY id LIMIT 1;
        END IF;

        UPDATE images SET user_id = admin_id WHERE user_id = 0;
    END
$$;

ALTER TABLE "images" ALTER COLUMN "user_id" DROP DEFAULT;

ALTER TABLE "images" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");