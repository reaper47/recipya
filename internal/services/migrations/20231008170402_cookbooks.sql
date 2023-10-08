-- +goose Up
CREATE TABLE cookbooks
(
    id      INTEGER PRIMARY KEY,
    title   TEXT NOT NULL,
    image   TEXT,
    count   INTEGER DEFAULT 0,
    user_id INTEGER REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE cookbook_recipes
(
    id          INTEGER PRIMARY KEY,
    cookbook_id INTEGER REFERENCES cookbooks (id) ON DELETE CASCADE,
    recipe_id   INTEGER REFERENCES recipes (id) ON DELETE CASCADE,
    UNIQUE (cookbook_id, recipe_id)
);

ALTER TABLE user_settings
    ADD COLUMN cookbooks_view INTEGER DEFAULT 0;

-- +goose StatementBegin
CREATE TRIGGER cookbook_recipes_insert
    AFTER INSERT
    ON cookbook_recipes
    FOR EACH ROW
BEGIN
    UPDATE cookbooks
    SET count = count + 1
    WHERE NEW.cookbook_id = id;
END;

CREATE TRIGGER cookbook_recipes_delete
    AFTER DELETE
    ON cookbook_recipes
    FOR EACH ROW
BEGIN
    UPDATE cookbooks
    SET count = count - 1
    WHERE OLD.cookbook_id = id;
END;
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER cookbook_recipes_insert;
DROP TRIGGER cookbook_recipes_delete;
DROP TABLE cookbook_recipes;
DROP TABLE cookbooks;
