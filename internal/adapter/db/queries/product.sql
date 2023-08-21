-- name: CreateProduct :one
INSERT INTO products (
    category_id, name, image, price, stock
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1
LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateProduct :one
UPDATE products
SET
    category_id = COALESCE(sqlc.narg(category_id), category_id),
    name = COALESCE(sqlc.narg(name), name),
    image = COALESCE(sqlc.narg(image), image),
    price = COALESCE(sqlc.narg(price), price),
    stock = COALESCE(sqlc.narg(stock), stock)
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;