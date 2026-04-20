-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    lesson_resources (
        id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        course_id INT UNSIGNED NOT NULL,
        FOREIGN KEY (course_id) REFERENCES course_details (id) ON DELETE CASCADE,
        lesson_id INT UNSIGNED NOT NULL,
        FOREIGN KEY (lesson_id) REFERENCES course_lessons (id) ON DELETE CASCADE,
        mime_type VARCHAR(50) NOT NULL,
        title VARCHAR(255) NOT NULL,
        file_path TEXT NOT NULL,
        file_size BIGINT NOT NULL,
        position INT DEFAULT 0,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS lesson_resources;

-- +goose StatementEnd