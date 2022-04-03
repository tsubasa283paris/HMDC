// Code generated by sqlc. DO NOT EDIT.
// source: users.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    id, password, name
) VALUES (
    $1, $2, $3
)
RETURNING id, password, name, created_at, updated_at, deleted_at
`

type CreateUserParams struct {
	ID       string `json:"id"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.ID, arg.Password, arg.Name)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Password,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
UPDATE users
SET deleted_at = NOW()
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUser = `-- name: GetUser :one
SELECT id, password, name, created_at, updated_at, deleted_at FROM users
WHERE id = $1 AND deleted_at IS NULL
`

func (q *Queries) GetUser(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Password,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const listUserDecks = `-- name: ListUserDecks :many
SELECT
    d.id,
    d.name,
    d.description,
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
WHERE d.user_id = $1
    AND d.deleted_at IS NULL
    AND NOT EXISTS (
        SELECT 1
        FROM league_decks l2
        WHERE l2.deck_id = l.deck_id
            AND l2.created_at > l.created_at
    )
`

type ListUserDecksRow struct {
	ID          int32         `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	LeagueID    sql.NullInt32 `json:"league_id"`
	NumDuel     int32         `json:"num_duel"`
	NumWin      int32         `json:"num_win"`
}

func (q *Queries) ListUserDecks(ctx context.Context, userID string) ([]ListUserDecksRow, error) {
	rows, err := q.db.QueryContext(ctx, listUserDecks, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListUserDecksRow
	for rows.Next() {
		var i ListUserDecksRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
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

const listUserDuelHistory = `-- name: ListUserDuelHistory :many
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
    dl1.created_at
FROM duels dl1
WHERE dl1.user_1_id = $1
    AND dl1.confirmed_at IS NOT NULL
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
    dl2.created_at
FROM duels dl2
WHERE dl2.user_2_id = $1
    AND dl2.confirmed_at IS NOT NULL
ORDER BY created_at
`

type ListUserDuelHistoryRow struct {
	ID             int32         `json:"id"`
	LeagueID       sql.NullInt32 `json:"league_id"`
	OpponentUserID string        `json:"opponent_user_id"`
	DeckID         int32         `json:"deck_id"`
	OpponentDeckID int32         `json:"opponent_deck_id"`
	Result         interface{}   `json:"result"`
	CreatedAt      time.Time     `json:"created_at"`
}

func (q *Queries) ListUserDuelHistory(ctx context.Context, user1ID string) ([]ListUserDuelHistoryRow, error) {
	rows, err := q.db.QueryContext(ctx, listUserDuelHistory, user1ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListUserDuelHistoryRow
	for rows.Next() {
		var i ListUserDuelHistoryRow
		if err := rows.Scan(
			&i.ID,
			&i.LeagueID,
			&i.OpponentUserID,
			&i.DeckID,
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

const listUserStats = `-- name: ListUserStats :many
SELECT
    l.id AS league_id,
    (SELECT (
        (
            SELECT COUNT(*)
            FROM duels dl
            WHERE dl.user_1_id = $1
                AND dl.league_id = l.id
                AND dl.confirmed_at IS NOT NULL
                AND dl.deleted_at IS NULL
        ) + (
            SELECT COUNT(*)
            FROM duels dl
            WHERE dl.user_2_id = $1
                AND dl.league_id = l.id
                AND dl.confirmed_at IS NOT NULL
                AND dl.deleted_at IS NULL
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
                AND dl.deleted_at IS NULL
        ) + (
            SELECT COUNT(*)
            FROM duels dl
            WHERE dl.user_2_id = $1
                AND dl.result = 2
                AND dl.league_id = l.id
                AND dl.confirmed_at IS NOT NULL
                AND dl.deleted_at IS NULL
        )
    )) AS num_win
FROM leagues l
`

type ListUserStatsRow struct {
	LeagueID int32 `json:"league_id"`
	NumDuel  int32 `json:"num_duel"`
	NumWin   int32 `json:"num_win"`
}

func (q *Queries) ListUserStats(ctx context.Context, user1ID string) ([]ListUserStatsRow, error) {
	rows, err := q.db.QueryContext(ctx, listUserStats, user1ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListUserStatsRow
	for rows.Next() {
		var i ListUserStatsRow
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

const listUsers = `-- name: ListUsers :many
SELECT 
    id,
    name
FROM users
WHERE deleted_at IS NULL
ORDER BY id
`

type ListUsersRow struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) ListUsers(ctx context.Context) ([]ListUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListUsersRow
	for rows.Next() {
		var i ListUsersRow
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
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

const updateUser = `-- name: UpdateUser :exec
UPDATE users
SET id = $2,
    password = $3,
    name = $4,
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
`

type UpdateUserParams struct {
	ID       string `json:"id"`
	ID_2     string `json:"id_2"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.ExecContext(ctx, updateUser,
		arg.ID,
		arg.ID_2,
		arg.Password,
		arg.Name,
	)
	return err
}
