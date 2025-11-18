-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
ADD COLUMN payment_method VARCHAR(256) NULL DEFAULT NULL,
ADD COLUMN transaction_id VARCHAR(256) NULL DEFAULT NULL;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders
DROP COLUMN payment_method,
DROP COLUMN transaction_id;

-- +goose StatementEnd