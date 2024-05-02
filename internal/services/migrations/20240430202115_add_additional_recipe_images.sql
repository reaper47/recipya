-- +goose Up
CREATE TABLE additional_images_recipe
(
    id        INTEGER PRIMARY KEY,
    recipe_id INTEGER NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    image     TEXT    NOT NULL,
    UNIQUE (recipe_id, image)
);

-- +goose Down
DROP TABLE additional_images_recipe;
