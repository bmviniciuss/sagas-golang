-- name: ListOrders :many
SELECT id, uuid, customer_id, amount, currency_code, status, created_at, updated_at
FROM orders.orders;

-- name: InsertOrder :exec
INSERT INTO orders.orders
	("uuid", customer_id, status, amount, currency_code, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;

-- name: GetOrder :one
SELECT id, uuid, customer_id, amount, currency_code, status, created_at, updated_at
FROM orders.orders WHERE uuid = $1;

-- name: UpdateOrder :exec
UPDATE orders.orders
SET status = $2, updated_at = $3
WHERE uuid = $1;
