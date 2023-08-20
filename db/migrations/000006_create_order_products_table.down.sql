ALTER TABLE
    IF EXISTS "order_products" DROP CONSTRAINT "fk_products_order_products";

ALTER TABLE
    IF EXISTS "order_products" DROP CONSTRAINT "fk_orders_order_products";

DROP TABLE IF EXISTS "order_products";