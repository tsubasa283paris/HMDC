// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"database/sql"
	"time"
)

type Deck struct {
	ID          int32        `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	UserID      string       `json:"user_id"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   sql.NullTime `json:"deleted_at"`
}

type Duel struct {
	ID          int32         `json:"id"`
	LeagueID    sql.NullInt32 `json:"league_id"`
	User1ID     string        `json:"user_1_id"`
	User2ID     string        `json:"user_2_id"`
	Deck1ID     int32         `json:"deck_1_id"`
	Deck2ID     int32         `json:"deck_2_id"`
	Result      int32         `json:"result"`
	CreatedAt   time.Time     `json:"created_at"`
	ConfirmedAt sql.NullTime  `json:"confirmed_at"`
	DeletedAt   sql.NullTime  `json:"deleted_at"`
}

type League struct {
	ID        int32        `json:"id"`
	Name      string       `json:"name"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type LeagueDeck struct {
	ID        int32     `json:"id"`
	LeagueID  int32     `json:"league_id"`
	DeckID    int32     `json:"deck_id"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID        string       `json:"id"`
	Password  string       `json:"password"`
	Name      string       `json:"name"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}
