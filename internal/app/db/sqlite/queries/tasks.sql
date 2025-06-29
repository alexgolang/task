-- name: CreateTask :exec
INSERT INTO tasks (id, title, description, status, priority, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetTask :one
SELECT * FROM tasks WHERE id = ?;

-- name: UpdateTask :exec
UPDATE tasks SET 
    title = COALESCE(NULLIF(sqlc.arg(title), ''), title),
    description = COALESCE(sqlc.narg(description), description),
    status = COALESCE(NULLIF(sqlc.arg(status), ''), status),
    priority = COALESCE(NULLIF(sqlc.arg(priority), ''), priority),
    updated_at = sqlc.arg(updated_at)
WHERE id = sqlc.arg(id);

-- name: DeleteTask :execrows
DELETE FROM tasks WHERE id = ?;

-- name: GetTasks :many
SELECT * FROM tasks;