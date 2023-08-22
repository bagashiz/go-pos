CREATE TABLE "orders" (
    "id" BIGSERIAL PRIMARY KEY,
    "user_id" bigint NOT NULL,
    "payment_id" bigint NOT NULL,
    "customer_name" varchar NOT NULL,
    "total_price" decimal(18, 2) NOT NULL,
    "total_paid" decimal(18, 2) NOT NULL,
    "total_return" decimal(18, 2) NOT NULL,
    "receipt_code" uuid NOT NULL DEFAULT (gen_random_uuid()),
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX "orders_customer_name" ON "orders" ("customer_name");

CREATE INDEX "orders_payment_id" ON "orders" ("payment_id");

CREATE INDEX "orders_user_id" ON "orders" ("user_id");

CREATE UNIQUE INDEX "receipt_code" ON "orders" ("receipt_code");

ALTER TABLE
    "orders"
ADD
    CONSTRAINT "fk_payments_orders" FOREIGN KEY ("payment_id") REFERENCES "payments" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE
    "orders"
ADD
    CONSTRAINT "fk_users_orders" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;