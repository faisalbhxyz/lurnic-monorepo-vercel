-- +goose Up
-- +goose StatementBegin
ALTER TABLE quiz_questions
ADD COLUMN correct_answer JSON NULL DEFAULT NULL;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE quiz_questions DROP COLUMN correct_answer;

-- +goose StatementEnd
