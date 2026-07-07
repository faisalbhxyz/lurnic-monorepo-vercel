-- +goose Up
ALTER TABLE course_certificate_settings
ADD COLUMN watermark_opacity TINYINT UNSIGNED NOT NULL DEFAULT 30;

ALTER TABLE student_certificates
ADD COLUMN watermark_opacity TINYINT UNSIGNED NOT NULL DEFAULT 30;

-- +goose Down
ALTER TABLE student_certificates DROP COLUMN watermark_opacity;

ALTER TABLE course_certificate_settings DROP COLUMN watermark_opacity;
