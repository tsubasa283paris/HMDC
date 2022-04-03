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

-- name: GetDeckStats :many
SELECT
    l.id AS league_id,
    (SELECT (
        (
            SELECT COUNT(*)
            FROM duels dl
            WHERE dl.deck_1_id = $1 AND dl.league_id = l.id
        ) + (
            SELECT COUNT(*)
            FROM duels dl
            WHERE dl.deck_2_id = $1 AND dl.league_id = l.id
        )
    )) AS num_duel,
    (SELECT (
        (
            SELECT COUNT(*)
            FROM duels dl
            WHERE dl.deck_1_id = $1 AND dl.result = 1 AND dl.league_id = l.id
        ) + (
            SELECT COUNT(*)
            FROM duels dl
            WHERE dl.deck_2_id = $1 AND dl.result = 2 AND dl.league_id = l.id
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
