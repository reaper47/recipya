-- +goose Up
ALTER TABLE tool_recipe
    ADD COLUMN quantity INTEGER DEFAULT 1;

ALTER TABLE tool_recipe
    ADD COLUMN tool_order INTEGER;

-- +goose Down
ALTER TABLE tool_recipe DROP COLUMN quantity;
ALTER TABLE tool_recipe DROP COLUMN tool_order;
