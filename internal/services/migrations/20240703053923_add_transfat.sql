-- +goose Up
ALTER TABLE nutrition
    ADD COLUMN trans_fat TEXT;

-- +goose Down
ALTER TABLE nutrition
    DROP COLUMN trans_fat;