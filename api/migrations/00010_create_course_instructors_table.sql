-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    course_instructors (
        id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        course_id INT UNSIGNED NOT NULL,
        FOREIGN KEY (course_id) REFERENCES course_details (id),
        instructor_id INT UNSIGNED NOT NULL,
        FOREIGN KEY (instructor_id) REFERENCES instructors (id),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS course_instructors;

-- +goose StatementEnd