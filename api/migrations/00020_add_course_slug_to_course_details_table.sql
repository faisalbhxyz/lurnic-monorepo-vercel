-- +goose Up
-- +goose StatementBegin
ALTER TABLE course_details
ADD COLUMN slug VARCHAR(256) NULL DEFAULT NULL;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE instructors
DROP COLUMN slug;

-- +goose StatementEnd