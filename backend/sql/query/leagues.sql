-- name: ListLeagues :many
SELECT 
    id,
    name
FROM leagues
WHERE deleted_at IS NULL
ORDER BY id;
