-- +goose Up
-- +goose StatementBegin
ALTER TABLE instructors
ADD COLUMN image VARCHAR(256) NULL DEFAULT NULL,
ADD COLUMN role VARCHAR(100) NULL DEFAULT NULL,
ADD COLUMN designation VARCHAR(100) NULL DEFAULT NULL;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE instructors
DROP COLUMN image,
DROP COLUMN role,
DROP COLUMN designation;

-- +goose StatementEnd