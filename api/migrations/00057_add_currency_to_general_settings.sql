-- +goose Up
ALTER TABLE general_settings
ADD COLUMN currency VARCHAR(10) NOT NULL DEFAULT 'BDT';

-- +goose Down
ALTER TABLE general_settings
DROP COLUMN currency;
