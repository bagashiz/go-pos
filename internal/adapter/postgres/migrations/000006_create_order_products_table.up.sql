CREATE TABLE "order_products" (
    "id" BIGSERIAL PRIMARY KEY,
    "order_id" bigint NOT NULL,
    "product_id" bigint NOT NULL,
    "quantity" bigint NOT NULL,
    "total_price" decimal(18, 2) NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX "order_product_order_id" ON "order_products" ("order_id");

CREATE INDEX "order_product_product_id" ON "order_products" ("product_id");

ALTER TABLE
    "order_products"
ADD
    CONSTRAINT "fk_orders_order_products" FOREIGN KEY ("order_id") REFERENCES "orders" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE
    "order_products"
ADD
    CONSTRAINT "fk_products_order_products" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;