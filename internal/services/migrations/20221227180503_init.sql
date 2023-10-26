-- +goose Up
CREATE TABLE users
(
    id              INTEGER PRIMARY KEY,
    email           TEXT      NOT NULL UNIQUE CHECK ( LOWER(email) = email ),
    hashed_password TEXT      NOT NULL,
    is_confirmed    INTEGER            DEFAULT 0,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE auth_tokens
(
    id             INTEGER PRIMARY KEY,
    selector       CHAR(12),
    hash_validator CHAR(64),
    expires        INTEGER DEFAULT (unixepoch('now', '+1 month')),
    user_id        INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE users;
DROP TABLE auth_tokens;
