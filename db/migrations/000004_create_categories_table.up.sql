CREATE TABLE "categories" (
    "id" BIGSERIAL PRIMARY KEY,
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "name" varchar NOT NULL
);

CREATE UNIQUE INDEX "category_name" ON "categories" ("name");