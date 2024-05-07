CREATE SCHEMA IF NOT EXISTS orders;

CREATE TABLE IF NOT EXISTS orders.orders (
  id serial PRIMARY KEY,
  uuid uuid NOT NULL UNIQUE DEFAULT gen_random_uuid(),
  global_id uuid NOT NULL UNIQUE,
  customer_id uuid NOT NULL,
  status VARCHAR(255) NOT NULL,
  amount bigint NOT NULL,
  currency_code VARCHAR(3) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_orders_global_id ON orders.orders (global_id);

CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON orders.orders (customer_id);


CREATE TABLE IF NOT EXISTS orders.order_items (
  id serial PRIMARY KEY,
  uuid uuid NOT NULL,
  quantity integer NOT NULL,
  unit_price bigint NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  order_id integer NOT NULL,
  CONSTRAINT fk_order_id FOREIGN KEY (order_id) REFERENCES orders.orders(id)
);
