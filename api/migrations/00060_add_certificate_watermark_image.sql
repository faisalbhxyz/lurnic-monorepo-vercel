-- +goose Up
ALTER TABLE course_certificate_settings
ADD COLUMN watermark_image VARCHAR(500) NULL;

ALTER TABLE student_certificates
ADD COLUMN watermark_image VARCHAR(500) NULL;

-- +goose Down
ALTER TABLE student_certificates DROP COLUMN watermark_image;

ALTER TABLE course_certificate_settings DROP COLUMN watermark_image;
