-- name: InsertTicket :one
INSERT INTO kitchen.ticket 
(uuid, customer_id, status, amount, currency_code, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;

-- name: InsertTicketItem :exec
INSERT INTO kitchen.ticket_items (uuid, quantity, unit_price, ticket_id)
VALUES ($1, $2, $3, $4);
