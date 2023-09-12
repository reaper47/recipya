-- +goose Up
CREATE TABLE share
(
    id         INTEGER PRIMARY KEY,
    link       TEXT      NOT NULL,
    user_id    INTEGER REFERENCES users (id) ON DELETE CASCADE,
    recipe_id  INTEGER REFERENCES recipes (id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (link, recipe_id)
);

-- +goose Down
DROP TABLE share;
