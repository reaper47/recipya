-- +goose Up
CREATE TABLE reports
(
    id                INTEGER PRIMARY KEY,
    report_type       INTEGER   NOT NULL REFERENCES report_types (id) ON DELETE CASCADE,
    created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    exec_time_ns INTEGER   NOT NULL,
    user_id           INTEGER   NOT NULL REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE report_types
(
    id   INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE report_logs
(
    id           INTEGER PRIMARY KEY,
    report_id    INTEGER NOT NULL REFERENCES reports (id) ON DELETE CASCADE,
    title        TEXT    NOT NULL,
    success      INTEGER NOT NULL,
    error_reason TEXT    NOT NULL
);

INSERT INTO report_types (name)
VALUES ('import');

-- +goose Down
DROP TABLE report_logs;
DROP TABLE reports;
DROP TABLE report_types;
