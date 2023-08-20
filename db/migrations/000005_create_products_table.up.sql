CREATE TABLE "products" (
    "id" BIGSERIAL PRIMARY KEY,
    "category_id" bigint NOT NULL,
    "sku" varchar NOT NULL,
    "name" varchar NOT NULL,
    "stock" bigint NOT NULL,
    "price" decimal(18, 2) NOT NULL,
    "image" varchar,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX "products_category_id" ON "products" ("category_id");

CREATE INDEX "products_name" ON "products" ("name");

CREATE UNIQUE INDEX "sku" ON "products" ("sku");

ALTER TABLE
    "products"
ADD
    CONSTRAINT "fk_categories_products" FOREIGN KEY ("category_id") REFERENCES "categories" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;