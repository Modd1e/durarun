
-- name: CreateJob :one
INSERT INTO jobs (queue, payload, status)
VALUES ($1, $2, 'available')
RETURNING *;
