-- name: GetDeck :one
SELECT
    d.id,
    d.name,
    d.description,
    d.user_id,
    l.league_id
FROM decks d
LEFT JOIN league_decks l ON d.id = l.deck_id
WHERE d.id = $1 AND d.deleted_at IS NULL
AND   NOT EXISTS (
    SELECT 1
    FROM league_decks l2
    WHERE l2.deck_id = l.deck_id
    AND   l2.created_at > l.created_at
);

-- name: UpdateDeck :exec
UPDATE decks
SET name = $2,
    description = $3,
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListDeckStats :many
SELECT
    l.id AS league_id,
    (SELECT (
        (
            SELECT COUNT(*)
            FROM duels dl
            WHERE dl.deck_1_id = $1
                AND dl.league_id = l.id
                AND dl.confirmed_at IS NOT NULL
                AND dl.deleted_at IS NULL
        ) + (
            SELECT COUNT(*)
            FROM duels dl
            WHERE dl.deck_2_id = $1
                AND dl.league_id = l.id
                AND dl.confirmed_at IS NOT NULL
                AND dl.deleted_at IS NULL
        )
    )) AS num_duel,
    (SELECT (
        (
            SELECT COUNT(*)
            FROM duels dl
            WHERE dl.deck_1_id = $1
                AND dl.result = 1
                AND dl.league_id = l.id
                AND dl.confirmed_at IS NOT NULL
                AND dl.deleted_at IS NULL
        ) + (
            SELECT COUNT(*)
            FROM duels dl
            WHERE dl.deck_2_id = $1
                AND dl.result = 2
                AND dl.league_id = l.id
                AND dl.confirmed_at IS NOT NULL
                AND dl.deleted_at IS NULL
        )
    )) AS num_win
FROM leagues l;

-- name: ListDecks :many
SELECT
    d.id,
    d.name,
    d.description,
    d.user_id,
    l.league_id,
    (SELECT (
        (
            SELECT COUNT(*)
            FROM duels dl
            WHERE dl.deck_1_id = d.id
                AND dl.confirmed_at IS NOT NULL
                AND dl.deleted_at IS NULL
        ) + (
            SELECT COUNT(*)
            FROM duels dl
            WHERE dl.deck_2_id = d.id
                AND dl.confirmed_at IS NOT NULL
                AND dl.deleted_at IS NULL
        )
    )) AS num_duel,
    (SELECT (
        (
            SELECT COUNT(*)
            FROM duels dl
            WHERE dl.deck_1_id = d.id
                AND dl.result = 1
                AND dl.confirmed_at IS NOT NULL
                AND dl.deleted_at IS NULL
        ) + (
            SELECT COUNT(*)
            FROM duels dl
            WHERE dl.deck_2_id = d.id
                AND dl.result = 2
                AND dl.confirmed_at IS NOT NULL
                AND dl.deleted_at IS NULL
        )
    )) AS num_win
FROM decks d
LEFT JOIN league_decks l ON d.id = l.deck_id
WHERE d.deleted_at IS NULL
AND   NOT EXISTS (
    SELECT 1
    FROM league_decks l2
    WHERE l2.deck_id = l.deck_id
    AND   l2.created_at > l.created_at
);

-- name: ListDeckDuelHistory :many
SELECT
    dl1.id,
    dl1.league_id,
    dl1.user_2_id AS opponent_user_id,
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
    dl1.created_at
FROM duels dl1
WHERE dl1.deck_1_id = $1
    AND dl1.confirmed_at IS NOT NULL
    AND dl1.deleted_at IS NULL
UNION
SELECT
    dl2.id,
    dl2.league_id,
    dl2.user_1_id AS opponent_user_id,
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
    dl2.created_at
FROM duels dl2
WHERE dl2.deck_2_id = $1
    AND dl2.confirmed_at IS NOT NULL
    AND dl2.deleted_at IS NULL
ORDER BY created_at;