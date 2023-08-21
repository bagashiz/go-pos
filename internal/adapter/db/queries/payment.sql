-- name: CreatePayment :one
INSERT INTO payments (
    name, type, logo
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetPayment :one
SELECT * FROM payments
WHERE id = $1
LIMIT 1;

-- name: ListPayments :many
SELECT * FROM payments
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdatePayment :one
UPDATE payments
SET
    name = COALESCE(sqlc.narg(name), name),
    type = COALESCE(sqlc.narg(type), type),
    logo = COALESCE(sqlc.narg(logo), logo)
WHERE id = $1
RETURNING *;

-- name: DeletePayment :exec
DELETE FROM payments
WHERE id = $1;