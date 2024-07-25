-- name: Create :one
INSERT INTO outbox
    (topic, event, headers)
VALUES
    ($1, $2, $3)
RETURNING id;

-- name: GetWaitList :many
SELECT id, topic, event, headers, created_at
FROM outbox
WHERE completed_at IS NULL;

-- name: SetComplete :exec
UPDATE outbox
SET completed_at = now()
WHERE id = $1;
