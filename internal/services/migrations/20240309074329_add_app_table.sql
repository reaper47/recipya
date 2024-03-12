-- +goose Up
CREATE TABLE app
(
    id                     INTEGER PRIMARY KEY,
    is_update_available    INTEGER            DEFAULT 0,
    updated_at             TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_last_checked_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementBegin
INSERT INTO app (id)
VALUES (1);

CREATE TRIGGER trig_app_update_check_auo
    AFTER UPDATE OF is_update_available
    ON app
    FOR EACH ROW
BEGIN
    UPDATE app
    SET update_last_checked_at = CURRENT_TIMESTAMP
    WHERE id = 1;
END;
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER trig_app_update_check_auo;
DROP TABLE app;
