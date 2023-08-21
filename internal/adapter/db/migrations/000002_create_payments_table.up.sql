CREATE TABLE "payments" (
    "id" BIGSERIAL PRIMARY KEY,
    "name" varchar NOT NULL,
    "type" varchar NOT NULL,
    "logo" varchar,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX "payment_name" ON "payments" ("name");