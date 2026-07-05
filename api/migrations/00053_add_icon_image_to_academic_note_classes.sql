-- +goose Up
-- +goose StatementBegin
ALTER TABLE academic_note_classes
ADD COLUMN icon_image TEXT NULL;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE academic_note_classes
DROP COLUMN icon_image;

-- +goose StatementEnd
