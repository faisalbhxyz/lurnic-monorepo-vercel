-- +goose Up
ALTER TABLE course_assignments ADD COLUMN position INT DEFAULT 0;

-- +goose Down
ALTER TABLE course_assignments DROP COLUMN position;
