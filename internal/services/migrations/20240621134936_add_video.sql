-- +goose Up
CREATE TABLE video_recipe
(
    id          INTEGER PRIMARY KEY,
    video       TEXT      NOT NULL,
    recipe_id   INTEGER   NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    duration    TEXT               DEFAULT '',
    content_url TEXT,
    embed_url TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (video, content_url, embed_url, recipe_id)
);

-- +goose Down
DROP TABLE video_recipe;
