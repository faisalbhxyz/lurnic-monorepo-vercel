-- +goose Up
-- TiDB does not support altering an existing column to add AUTO_INCREMENT (Error 8200).
-- Use the one-off fixer `go run ./cmd/fixautoinc` which performs a shadow-table swap safely.
-- +goose StatementBegin
SELECT 1;
-- +goose StatementEnd

-- +goose Down
-- (no-op) keeping AUTO_INCREMENT is safe
-- +goose StatementBegin
SELECT 1;
-- +goose StatementEnd

