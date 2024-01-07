-- +goose Up
ALTER TABLE nutrition
    ADD COLUMN is_per_serving INTEGER DEFAULT 0;

-- +goose Down
ALTER TABLE nutrition
    DROP COLUMN is_per_serving;

