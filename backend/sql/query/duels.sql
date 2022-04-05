-- name: ListDuelsWithLimit :many
SELECT
    id,
    league_id,
    user_1_id,
    user_2_id,
    deck_1_id,
    deck_2_id,
    result,
    created_at
FROM duels
WHERE confirmed_at IS NOT NULL
    AND deleted_at IS NULL
ORDER BY created_at
LIMIT $1;

-- name: ListUnconfirmedDuelsByUser :many
SELECT
    dl1.id,
    dl1.league_id,
    dl1.user_2_id AS opponent_user_id,
    dl1.deck_1_id AS deck_id,
    dl1.deck_2_id AS opponent_deck_id,
    (SELECT
        CASE
            WHEN dl1.result = -1 THEN '-'
            WHEN dl1.result = 0 THEN 'draw'
            WHEN dl1.result = 1 THEN 'win'
            WHEN dl1.result = 2 THEN 'lose'
            ELSE 'undefined'
        END
    ) AS result,
    dl1.created_at,
    dl1.created_by
FROM duels dl1
WHERE dl1.user_1_id = $1
    AND dl1.confirmed_at IS NULL
    AND dl1.deleted_at IS NULL
    AND (dl1.created_by = 1 OR dl1.created_by = 2)
UNION
SELECT
    dl2.id,
    dl2.league_id,
    dl2.user_1_id AS opponent_user_id,
    dl2.deck_2_id AS deck_id,
    dl2.deck_1_id AS opponent_deck_id,
    (SELECT
        CASE
            WHEN dl2.result = -1 THEN '-'
            WHEN dl2.result = 0 THEN 'draw'
            WHEN dl2.result = 1 THEN 'lose'
            WHEN dl2.result = 2 THEN 'win'
            ELSE 'undefined'
        END
    ) AS result,
    dl2.created_at,
    dl2.created_by
FROM duels dl2
WHERE dl2.user_2_id = $1
    AND dl2.confirmed_at IS NULL
    AND dl2.deleted_at IS NULL
    AND (dl2.created_by = 1 OR dl2.created_by = 2)
ORDER BY created_at;

-- name: CreateUnconfirmedDuel :exec
INSERT INTO duels (
    league_id,
    user_1_id,
    user_2_id,
    deck_1_id,
    deck_2_id,
    result,
    created_by
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
);

-- name: GetDuel :one
SELECT *
FROM duels
WHERE id = $1
    AND deleted_at IS NULL;

-- name: ConfirmDuel :exec
UPDATE duels
SET confirmed_at = NOW()
WHERE id = $1
    AND deleted_at IS NULL;

-- name: DeleteDuel :exec
UPDATE duels
SET deleted_at = NOW()
WHERE id = $1
    AND deleted_at IS NULL;
