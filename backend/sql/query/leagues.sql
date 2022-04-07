-- name: GetLeague :one
SELECT
    *
FROM leagues
WHERE id = $1;

-- name: ListLeagues :many
SELECT 
    id,
    name
FROM leagues
WHERE deleted_at IS NULL
ORDER BY id;
