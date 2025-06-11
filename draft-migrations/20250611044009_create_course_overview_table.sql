-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    course_overviews (
        id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        course_id INT UNSIGNED NOT NULL,
        outcomes TEXT NOT NULL,
        target_audience TEXT NULL,
        duration_hours INT NULL,
        duration_mins INT NULL,
        materials_included TEXT NULL,
        requirements TEXT NULL,
        CONSTRAINT fk_course_overview FOREIGN KEY (course_id) REFERENCES courses (id)
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS course_overviews;

-- +goose StatementEnd