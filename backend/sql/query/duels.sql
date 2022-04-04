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
