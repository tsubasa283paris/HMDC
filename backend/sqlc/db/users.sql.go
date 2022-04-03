// Code generated by sqlc. DO NOT EDIT.
// source: users.sql

package db

import (
	"context"
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

const getUserStats = `-- name: GetUserStats :many
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
FROM leagues l
`

type GetUserStatsRow struct {
	LeagueID int32 `json:"league_id"`
	NumDuel  int32 `json:"num_duel"`
	NumWin   int32 `json:"num_win"`
}

func (q *Queries) GetUserStats(ctx context.Context, user1ID string) ([]GetUserStatsRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserStats, user1ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserStatsRow
	for rows.Next() {
		var i GetUserStatsRow
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
