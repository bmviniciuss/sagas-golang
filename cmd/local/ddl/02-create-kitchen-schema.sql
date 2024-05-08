CREATE SCHEMA IF NOT EXISTS kitchen;

CREATE TABLE IF NOT EXISTS kitchen.ticket (
  id serial PRIMARY KEY,
  uuid uuid NOT NULL UNIQUE DEFAULT gen_random_uuid(),
  customer_id uuid NOT NULL,
  status VARCHAR(255) NOT NULL,
  amount bigint NOT NULL,
  currency_code VARCHAR(3) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_kitchen_uuid ON kitchen.ticket (uuid);

CREATE INDEX IF NOT EXISTS idx_kitchen_customer_id ON kitchen.ticket (customer_id);

CREATE TABLE IF NOT EXISTS kitchen.ticket_items (
  id serial PRIMARY KEY,
  uuid uuid NOT NULL UNIQUE DEFAULT gen_random_uuid(),
  quantity int NOT NULL,
  unit_price bigint NOT NULL,
  ticket_id int NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT fk_ticket_id FOREIGN KEY (ticket_id) REFERENCES kitchen.ticket(id)
);

CREATE INDEX IF NOT EXISTS idx_kitchen_ticket_id ON kitchen.ticket_items (ticket_id);
