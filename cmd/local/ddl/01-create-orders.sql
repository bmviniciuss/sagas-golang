CREATE SCHEMA IF NOT EXISTS orders;

CREATE TABLE IF NOT EXISTS orders.orders (
  id serial PRIMARY KEY,
  uuid uuid NOT NULL UNIQUE DEFAULT gen_random_uuid(),
  customer_id uuid NOT NULL,
  status VARCHAR(255) NOT NULL,
  amount bigint NOT NULL,
  currency_code VARCHAR(3) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_orders_uuid ON orders.orders (uuid);

CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON orders.orders (customer_id);
