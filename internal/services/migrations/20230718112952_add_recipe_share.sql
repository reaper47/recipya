-- +goose Up
CREATE TABLE share
(
    id         INTEGER PRIMARY KEY,
    link       TEXT      NOT NULL,
    recipe_id  INTEGER REFERENCES recipes (id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (link, recipe_id)
);

-- +goose Down
DROP TABLE share;
