-- +goose Up
ALTER TABLE quiz_submissions
    ADD COLUMN instructor_feedback TEXT NULL DEFAULT NULL;

-- +goose Down
ALTER TABLE quiz_submissions
    DROP COLUMN instructor_feedback;
