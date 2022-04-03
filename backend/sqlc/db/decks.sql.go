// Code generated by sqlc. DO NOT EDIT.
// source: decks.sql

package db

import (
	"context"
	"database/sql"
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

const getDeckStats = `-- name: GetDeckStats :many
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
FROM leagues l
`

type GetDeckStatsRow struct {
	LeagueID int32 `json:"league_id"`
	NumDuel  int32 `json:"num_duel"`
	NumWin   int32 `json:"num_win"`
}

func (q *Queries) GetDeckStats(ctx context.Context, deck1ID int32) ([]GetDeckStatsRow, error) {
	rows, err := q.db.QueryContext(ctx, getDeckStats, deck1ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetDeckStatsRow
	for rows.Next() {
		var i GetDeckStatsRow
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
