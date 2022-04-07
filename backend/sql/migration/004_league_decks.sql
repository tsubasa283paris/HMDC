-- +migrate Up
CREATE TABLE league_decks (
    id INTEGER GENERATED ALWAYS AS IDENTITY,
    league_id INTEGER NOT NULL,
    deck_id INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY(id),
    CONSTRAINT fk_league
        FOREIGN KEY(league_id)
            REFERENCES leagues(id)
            ON DELETE SET NULL,
    CONSTRAINT fk_deck
        FOREIGN KEY(deck_id)
            REFERENCES decks(id)
            ON DELETE SET NULL
);


-- +migrate Down
DROP TABLE league_decks;