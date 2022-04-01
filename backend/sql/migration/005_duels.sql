-- +migrate Up
CREATE TABLE duels (
    id INTEGER GENERATED ALWAYS AS IDENTITY,
    league_id INTEGER NOT NULL,
    user_1_id VARCHAR(255) NOT NULL,
    user_2_id VARCHAR(255) NOT NULL,
    deck_1_id INTEGER NOT NULL,
    deck_2_id INTEGER NOT NULL,
    result INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    confirmed_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    PRIMARY KEY(id),
    CONSTRAINT fk_league
        FOREIGN KEY(league_id)
            REFERENCES leagues(id)
            ON DELETE SET NULL,
    CONSTRAINT fk_user_1
        FOREIGN KEY(user_1_id)
            REFERENCES users(id)
            ON DELETE SET NULL,
    CONSTRAINT fk_user_2
        FOREIGN KEY(user_2_id)
            REFERENCES users(id)
            ON DELETE SET NULL,
    CONSTRAINT fk_deck_1
        FOREIGN KEY(deck_1_id)
            REFERENCES decks(id)
            ON DELETE SET NULL,
    CONSTRAINT fk_deck_2
        FOREIGN KEY(deck_2_id)
            REFERENCES decks(id)
            ON DELETE SET NULL
);


-- +migrate Down
DROP TABLE duels;