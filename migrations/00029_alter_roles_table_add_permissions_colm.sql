-- +goose Up
-- +goose StatementBegin
ALTER TABLE `roles`
ADD COLUMN `permissions` JSON NULL DEFAULT NULL;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE `roles`
DROP COLUMN `permissions`;

-- +goose StatementEnd