
-- name: CreateJob :one
INSERT INTO jobs (queue, payload, status)
VALUES ($1, $2, 'available')
RETURNING *;

-- name: CountJobs :one
SELECT COUNT(*) FROM jobs;
