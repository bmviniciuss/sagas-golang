CREATE SCHEMA IF NOT EXISTS orders;

CREATE TABLE IF NOT EXISTS orders.orders (
  id serial PRIMARY KEY,
  global_id uuid NOT NULL UNIQUE,
  uuid uuid NOT NULL UNIQUE DEFAULT gen_random_uuid(),
  client_id uuid NOT NULL,
  customer_id uuid NOT NULL,
  "date" timestamptz NOT NULL DEFAULT now(),
  status VARCHAR(255) NOT NULL,
  total bigint NOT NULL,
  currency_code VARCHAR(3) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_orders_client_id ON orders.orders (client_id);
CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON orders.orders (customer_id);
CREATE INDEX IF NOT EXISTS idx_orders_date ON orders.orders ("date");
