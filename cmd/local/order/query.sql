-- name: ListOrders :many
SELECT id, uuid, customer_id, amount, currency_code, status, created_at, updated_at
FROM orders.orders;

-- name: InsertOrder :one
INSERT INTO orders.orders
	("uuid", customer_id, status, amount, currency_code, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;

-- name: InsertOrderItem :exec
INSERT INTO orders.order_items
("uuid", quantity, unit_price, order_id, created_at, updated_at)
VALUES($1, $2, $3, $4, now(), now());

-- name: GetOrder :one
SELECT id, uuid, customer_id, amount, currency_code, status, created_at, updated_at
FROM orders.orders WHERE uuid = $1;

-- name: GetOrderItems :many
SELECT oi.* FROM orders.order_items oi
join orders.orders o on o.id  = oi.order_id
where o.uuid = $1;
