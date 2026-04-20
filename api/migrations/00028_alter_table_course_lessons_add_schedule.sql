-- +goose Up
-- +goose StatementBegin
ALTER TABLE course_lessons
ADD COLUMN is_scheduled BOOLEAN NOT NULL DEFAULT FALSE,
ADD COLUMN schedule_date DATE NULL,
ADD COLUMN schedule_time TIME NULL,
ADD COLUMN show_comming_soon BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE course_lessons
DROP COLUMN is_scheduled,
DROP COLUMN schedule_date,
DROP COLUMN schedule_time,
DROP COLUMN show_comming_soon;

-- +goose StatementEnd