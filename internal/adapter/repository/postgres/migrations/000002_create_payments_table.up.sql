CREATE TYPE "payments_type_enum" AS ENUM ('CASH', 'E-WALLET', 'EDC');

CREATE TABLE "payments" (
    "id" BIGSERIAL PRIMARY KEY,
    "name" varchar NOT NULL,
    "type" payments_type_enum NOT NULL,
    "logo" varchar,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX "payment_name" ON "payments" ("name");