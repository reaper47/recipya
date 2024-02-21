-- +goose Up
-- +goose StatementBegin
CREATE TRIGGER recipe_delete
    AFTER DELETE
    ON recipes
    FOR EACH ROW
BEGIN
    DELETE
    FROM recipes_fts
    WHERE id = OLD.id;
END;
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER recipe_delete;
