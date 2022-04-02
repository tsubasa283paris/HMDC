-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListUsers :many
SELECT 
    id,
    name
FROM users
WHERE deleted_at IS NULL
ORDER BY id;

-- name: CreateUser :one
INSERT INTO users (
    id, password, name
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: UpdateUser :exec
UPDATE users
SET id = $2,
    password = $3,
    name = $4,
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: DeleteUser :exec
UPDATE users
SET deleted_at = NOW()
WHERE id = $1;