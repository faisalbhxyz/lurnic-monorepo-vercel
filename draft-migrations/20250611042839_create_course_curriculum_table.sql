-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    course_curriculums (
        id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        course_id INT UNSIGNED NOT NULL,
        position INTEGER DEFAULT 0,
        title TEXT NOT NULL,
        description TEXT NULL,
        CONSTRAINT fk_course FOREIGN KEY (course_id) REFERENCES courses (id)
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS course_curriculums;

-- +goose StatementEnd