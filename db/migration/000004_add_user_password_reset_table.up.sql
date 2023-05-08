CREATE TABLE "user_password_reset" (
                                       "user_id" int NOT NULL,
                                       "code" varchar NOT NULL,
                                       "expired_at" timestamptz NOT NULL DEFAULT (now() + INTERVAL '15 minutes')
);


ALTER TABLE "user_password_reset" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");