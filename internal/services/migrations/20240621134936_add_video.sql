-- +goose Up
CREATE TABLE video_recipe
(
    id        INTEGER PRIMARY KEY,
    video     TEXT    NOT NULL,
    recipe_id INTEGER NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    UNIQUE (video, recipe_id)
);

-- +goose Down
DROP TABLE video_recipe;
