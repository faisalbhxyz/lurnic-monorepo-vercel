-- +goose Up
-- +goose StatementBegin
ALTER TABLE course_general_settings
ADD COLUMN IF NOT EXISTS sub_category_id INT UNSIGNED NULL DEFAULT NULL;

-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE course_general_settings
ADD CONSTRAINT fk_sub_category FOREIGN KEY (sub_category_id) REFERENCES sub_categories (id) ON DELETE CASCADE;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE course_general_settings
DROP FOREIGN KEY fk_sub_category;
DROP COLUMN sub_category_id;

-- +goose StatementEnd