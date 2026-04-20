-- +goose Up
-- +goose StatementBegin
ALTER TABLE course_lessons
DROP COLUMN resources;

-- +goose StatementEnd