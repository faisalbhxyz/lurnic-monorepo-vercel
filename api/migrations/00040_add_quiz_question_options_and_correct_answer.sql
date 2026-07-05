-- +goose Up
-- +goose StatementBegin
ALTER TABLE quiz_questions
ADD COLUMN options JSON NULL DEFAULT NULL;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE quiz_questions DROP COLUMN options;

-- +goose StatementEnd
