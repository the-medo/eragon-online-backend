CREATE TABLE "users" (
  "id" integer PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "email" varchar NOT NULL,
  "avatar" varchar,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "roles" (
  "id" integer PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL
);

CREATE TABLE "user_roles" (
  "user_id" integer NOT NULL,
  "role_id" integer NOT NULL
);

CREATE TABLE "chat" (
  "id" integer PRIMARY KEY,
  "user_id" integer NOT NULL,
  "text" varchar NOT NULL,
  "created_at" timestamp DEFAULT (now())
);

ALTER TABLE "chat" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_roles" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_roles" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");
