ALTER TABLE
    IF EXISTS "orders" DROP CONSTRAINT "fk_users_orders";

ALTER TABLE
    IF EXISTS "orders" DROP CONSTRAINT "fk_payments_orders";

DROP TABLE IF EXISTS "orders";