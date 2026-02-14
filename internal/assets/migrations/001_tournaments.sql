-- +migrate Up

CREATE TABLE tournaments(
    id TEXT PRIMARY KEY,
    name TEXT,
    format TEXT NOT NULL,
    lan BOOLEAN DEFAULT FALSE,
    region TEXT NOT NULL,
    prize_pool NUMERIC,
    prize_distribution JSONB,
    date TIMESTAMP NOT NULL
);

-- +migrate Down
DROP TABLE tournaments;