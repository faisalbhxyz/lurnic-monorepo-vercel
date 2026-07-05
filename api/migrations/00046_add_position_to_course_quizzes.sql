-- +goose Up
ALTER TABLE course_quizzes ADD COLUMN position INT DEFAULT 0;

-- +goose Down
ALTER TABLE course_quizzes DROP COLUMN position;
