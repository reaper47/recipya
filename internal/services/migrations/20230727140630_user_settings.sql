-- +goose Up
CREATE TABLE measurement_systems
(
    id   INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE user_settings
(
    id                    INTEGER PRIMARY KEY,
    user_id               INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    measurement_system_id INTEGER REFERENCES measurement_systems (id) ON DELETE CASCADE DEFAULT 1,
    convert_automatically INTEGER NOT NULL DEFAULT 0
);

INSERT INTO measurement_systems (name)
VALUES ('imperial'),
       ('metric');

INSERT INTO user_settings (user_id, measurement_system_id)
SELECT id, 1
FROM users;

-- +goose StatementBegin
CREATE TRIGGER user_settings_default
    AFTER INSERT
    ON users
    FOR EACH ROW
BEGIN
    INSERT INTO user_settings (user_id)
    VALUES (NEW.id);
END;
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER user_settings_default;
DROP TABLE measurement_systems;
DROP TABLE user_settings;
