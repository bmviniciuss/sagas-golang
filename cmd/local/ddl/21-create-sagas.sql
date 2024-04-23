CREATE SCHEMA IF NOT EXISTS sagas;

CREATE TABLE IF NOT EXISTS sagas.workflows (
  id uuid PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE,
  state JSONB NOT NULL DEFAULT '{}',
  reply_channel VARCHAR(255) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_sagas_workflows_id ON sagas.workflows (id);
CREATE INDEX IF NOT EXISTS idx_sagas_workflows_name ON sagas.workflows (name);

CREATE TABLE IF NOT EXISTS sagas.steps (
  id uuid PRIMARY KEY,
  workflow_id uuid NOT NULL,
  order_number smallint NOT NULL,
  name VARCHAR(255) NOT NULL UNIQUE,
  service_name VARCHAR(255) NOT NULL,
  compensable BOOLEAN NOT NULL DEFAULT FALSE,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  FOREIGN KEY (workflow_id) REFERENCES sagas.workflows (id)
);

CREATE INDEX IF NOT EXISTS idx_sagas_steps_id ON sagas.steps (id);
CREATE INDEX IF NOT EXISTS idx_sagas_steps_workflow_id ON sagas.steps (workflow_id);
CREATE INDEX IF NOT EXISTS idx_sagas_steps_name ON sagas.steps (name);


CREATE TABLE IF NOT EXISTS sagas.executions (
  id uuid PRIMARY KEY,
  workflow_id uuid NOT NULL,
  state JSONB NOT NULL DEFAULT '{}',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  FOREIGN KEY (workflow_id) REFERENCES sagas.workflows (id)
);

CREATE INDEX IF NOT EXISTS idx_sagas_executions_id ON sagas.executions (id);
CREATE INDEX IF NOT EXISTS idx_sagas_executions_workflow_id ON sagas.executions (workflow_id);

INSERT INTO sagas.workflows (id, name, reply_channel) 
VALUES ('2ef23373-9c01-4603-be2f-8e80552eb9a4', 'create_order', 'saga.create-order.v1.response');

INSERT INTO sagas.steps (id, workflow_id, order_number, name, service_name, compensable)
VALUES 
  ('4a4578ff-3602-4ad0-b262-6827c6ebc985', '2ef23373-9c01-4603-be2f-8e80552eb9a4', 1, 'create_order', 'order', true),
  ('22d7c4bb-e751-4b47-a7a0-903ee5d3996e', '2ef23373-9c01-4603-be2f-8e80552eb9a4', 2, 'create_payment', 'customer', false);


