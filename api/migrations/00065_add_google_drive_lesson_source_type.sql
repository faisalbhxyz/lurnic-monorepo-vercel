-- +goose Up
-- +goose StatementBegin
ALTER TABLE course_lessons
MODIFY COLUMN source_type ENUM(
    'youtube',
    'vimeo',
    'sound_cloud',
    'spotify',
    'custom_code',
    'upload',
    'google_drive'
) NOT NULL DEFAULT 'upload';

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE course_lessons
MODIFY COLUMN source_type ENUM(
    'youtube',
    'vimeo',
    'sound_cloud',
    'spotify',
    'custom_code',
    'upload'
) NOT NULL DEFAULT 'youtube';

-- +goose StatementEnd
