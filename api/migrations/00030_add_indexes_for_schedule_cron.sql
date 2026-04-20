-- +goose Up
-- +goose StatementBegin
-- Speed up cron queries: WHERE is_scheduled = true (course publish)
CREATE INDEX idx_course_details_is_scheduled ON course_details (is_scheduled);

-- Speed up cron: WHERE is_scheduled = true AND is_published = false (lesson publish)
CREATE INDEX idx_course_lessons_scheduled_published ON course_lessons (is_scheduled, is_published);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_course_lessons_scheduled_published ON course_lessons;
DROP INDEX idx_course_details_is_scheduled ON course_details;

-- +goose StatementEnd
