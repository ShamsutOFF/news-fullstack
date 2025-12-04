-- +goose Up
CREATE TABLE sessions (
    key VARCHAR(64) PRIMARY KEY,
    data BYTEA NOT NULL,
    exp TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_sessions_exp ON sessions (exp);

-- +goose Down
DROP INDEX IF EXISTS idx_sessions_exp;
DROP TABLE IF EXISTS sessions;