-- +goose Up
-- +goose StatementBegin
CREATE TRIGGER users_delete
    AFTER DELETE
    ON users
    FOR EACH ROW
BEGIN
    DELETE
    FROM recipes
    WHERE id IN (SELECT id
                 FROM recipes_fts AS r
                 WHERE r.user_id = OLD.id);

    DELETE
    FROM recipes_fts
    WHERE user_id = OLD.id;

    DELETE
    FROM cookbooks
    WHERE user_id = OLD.id;

    DELETE
    FROM cookbooks_fts
    WHERE user_id = OLD.id;
END;
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER users_delete;
