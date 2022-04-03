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

-- name: GetUserStats :many
SELECT
    l.id AS league_id,
    (SELECT (
        (
            SELECT COUNT(*)
            FROM duels dl
            WHERE dl.user_1_id = $1
                AND dl.league_id = l.id
                AND dl.confirmed_at IS NOT NULL
        ) + (
            SELECT COUNT(*)
            FROM duels dl
            WHERE dl.user_2_id = $1
                AND dl.league_id = l.id
                AND dl.confirmed_at IS NOT NULL
        )
    )) AS num_duel,
    (SELECT (
        (
            SELECT COUNT(*)
            FROM duels dl
            WHERE dl.user_1_id = $1
                AND dl.result = 1
                AND dl.league_id = l.id
                AND dl.confirmed_at IS NOT NULL
        ) + (
            SELECT COUNT(*)
            FROM duels dl
            WHERE dl.user_2_id = $1
                AND dl.result = 2
                AND dl.league_id = l.id
                AND dl.confirmed_at IS NOT NULL
        )
    )) AS num_win
FROM leagues l;
