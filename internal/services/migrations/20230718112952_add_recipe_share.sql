-- +goose Up
CREATE TABLE share_cookbooks
(
    id         INTEGER PRIMARY KEY,
    link       TEXT      NOT NULL,
    user_id    INTEGER REFERENCES users (id) ON DELETE CASCADE,
    cookbook_id INTEGER REFERENCES cookbooks (id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (link, cookbook_id)
);

CREATE TABLE share_recipes
(
    id         INTEGER PRIMARY KEY,
    link       TEXT      NOT NULL,
    user_id    INTEGER REFERENCES users (id) ON DELETE CASCADE,
    recipe_id  INTEGER REFERENCES recipes (id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (link, recipe_id)
);

-- +goose Down
DROP TABLE share_cookbooks;
DROP TABLE share_recipes;
