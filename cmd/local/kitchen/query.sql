-- name: InsertTicket :one
INSERT INTO kitchen.ticket 
(uuid, customer_id, status, amount, currency_code, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;

-- name: InsertTicketItem :exec
INSERT INTO kitchen.ticket_items (uuid, quantity, unit_price, ticket_id)
VALUES ($1, $2, $3, $4);

-- name: UpdateTicketStatus :exec
UPDATE kitchen.ticket SET status = $1, updated_at = $2 WHERE uuid = $3;

-- name: GetTicket :one
SELECT id, uuid, customer_id, status, amount, currency_code, created_at, updated_at FROM kitchen.ticket WHERE uuid = $1;

-- name: GetTicketItems :many
SELECT uuid, quantity, unit_price FROM kitchen.ticket_items WHERE ticket_id = $1;
