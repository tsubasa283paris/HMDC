// Code generated by sqlc. DO NOT EDIT.
// source: decks.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const getDeck = `-- name: GetDeck :one
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
)
`

type GetDeckRow struct {
	ID          int32         `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	UserID      string        `json:"user_id"`
	LeagueID    sql.NullInt32 `json:"league_id"`
}

func (q *Queries) GetDeck(ctx context.Context, id int32) (GetDeckRow, error) {
	row := q.db.QueryRowContext(ctx, getDeck, id)
	var i GetDeckRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.UserID,
		&i.LeagueID,
	)
	return i, err
}

const listDeckDuelHistory = `-- name: ListDeckDuelHistory :many
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
ORDER BY created_at
`

type ListDeckDuelHistoryRow struct {
	ID             int32         `json:"id"`
	LeagueID       sql.NullInt32 `json:"league_id"`
	OpponentUserID string        `json:"opponent_user_id"`
	OpponentDeckID int32         `json:"opponent_deck_id"`
	Result         interface{}   `json:"result"`
	CreatedAt      time.Time     `json:"created_at"`
}

func (q *Queries) ListDeckDuelHistory(ctx context.Context, deck1ID int32) ([]ListDeckDuelHistoryRow, error) {
	rows, err := q.db.QueryContext(ctx, listDeckDuelHistory, deck1ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListDeckDuelHistoryRow
	for rows.Next() {
		var i ListDeckDuelHistoryRow
		if err := rows.Scan(
			&i.ID,
			&i.LeagueID,
			&i.OpponentUserID,
			&i.OpponentDeckID,
			&i.Result,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listDeckStats = `-- name: ListDeckStats :many
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
FROM leagues l
`

type ListDeckStatsRow struct {
	LeagueID int32 `json:"league_id"`
	NumDuel  int32 `json:"num_duel"`
	NumWin   int32 `json:"num_win"`
}

func (q *Queries) ListDeckStats(ctx context.Context, deck1ID int32) ([]ListDeckStatsRow, error) {
	rows, err := q.db.QueryContext(ctx, listDeckStats, deck1ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListDeckStatsRow
	for rows.Next() {
		var i ListDeckStatsRow
		if err := rows.Scan(&i.LeagueID, &i.NumDuel, &i.NumWin); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listDecks = `-- name: ListDecks :many
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
)
`

type ListDecksRow struct {
	ID          int32         `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	UserID      string        `json:"user_id"`
	LeagueID    sql.NullInt32 `json:"league_id"`
	NumDuel     int32         `json:"num_duel"`
	NumWin      int32         `json:"num_win"`
}

func (q *Queries) ListDecks(ctx context.Context) ([]ListDecksRow, error) {
	rows, err := q.db.QueryContext(ctx, listDecks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListDecksRow
	for rows.Next() {
		var i ListDecksRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.UserID,
			&i.LeagueID,
			&i.NumDuel,
			&i.NumWin,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateDeck = `-- name: UpdateDeck :exec
UPDATE decks
SET name = $2,
    description = $3,
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
`

type UpdateDeckParams struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (q *Queries) UpdateDeck(ctx context.Context, arg UpdateDeckParams) error {
	_, err := q.db.ExecContext(ctx, updateDeck, arg.ID, arg.Name, arg.Description)
	return err
}
