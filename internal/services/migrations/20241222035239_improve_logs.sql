-- +goose Up
ALTER TABLE report_logs ADD COLUMN warning INTEGER NOT NULL DEFAULT 0;
ALTER TABLE report_logs ADD COLUMN action TEXT;

-- +goose Down
ALTER TABLE report_logs DROP COLUMN warning;
ALTER TABLE report_logs DROP COLUMN action;
