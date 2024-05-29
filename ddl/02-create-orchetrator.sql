CREATE SCHEMA IF NOT EXISTS sagas;

CREATE TABLE IF NOT EXISTS sagas.executions (
  id serial PRIMARY KEY,
  uuid uuid NOT NULL UNIQUE DEFAULT gen_random_uuid(),
  workflow_name varchar(255) NOT NULL,
  state jsonb NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);
