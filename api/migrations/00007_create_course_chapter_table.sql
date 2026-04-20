-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    course_chapters (
        id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        position INT DEFAULT 0,
        title VARCHAR(255) NOT NULL,
        description TEXT NULL,
        access ENUM ('draft', 'published') NOT NULL DEFAULT 'published',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        course_id INT UNSIGNED NOT NULL,
        FOREIGN KEY (course_id) REFERENCES course_details (id) ON DELETE CASCADE
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS course_chapters;

-- +goose StatementEnd