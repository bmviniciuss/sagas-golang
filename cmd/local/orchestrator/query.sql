-- name: InsertExecution :exec
INSERT INTO sagas.executions
	("uuid", workflow_name, state, created_at, updated_at)
VALUES ($1, $2, $3, now(), now()) RETURNING id;

-- name: UpdateExecution :exec
UPDATE sagas.executions
SET state = $2, updated_at = now()
WHERE uuid = $1;

-- name: FindExecutionByUUID :one
SELECT id, uuid, workflow_name, state, created_at, updated_at
FROM sagas.executions
WHERE uuid = $1 LIMIT 1;
