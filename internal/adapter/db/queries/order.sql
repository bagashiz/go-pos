-- name: CreateOrder :one
INSERT INTO orders (
    user_id,
    payment_id,
    customer_name,
    total_price,
    total_paid,
    total_return
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders
WHERE id = $1
LIMIT 1;

-- name: ListOrders :many
SELECT * FROM orders
ORDER BY id
LIMIT $1
OFFSET $2;
