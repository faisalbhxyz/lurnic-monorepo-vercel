-- +goose Up
-- +goose StatementBegin
ALTER TABLE course_details
ADD COLUMN position INT DEFAULT 0;

-- +goose StatementEnd