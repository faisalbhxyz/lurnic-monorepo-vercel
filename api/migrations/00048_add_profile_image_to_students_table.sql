-- +goose Up
ALTER TABLE students
ADD COLUMN profile_image VARCHAR(256) NULL DEFAULT NULL;

-- +goose Down
ALTER TABLE students
DROP COLUMN profile_image;
