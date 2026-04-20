-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    course_general_settings (
        id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        difficulty_level ENUM ('all', 'beginner', 'intermediate', 'expert') DEFAULT 'all',
        maximum_student INTEGER DEFAULT 0,
        language VARCHAR(100) DEFAULT 'english',
        duration VARCHAR(255) NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        category_id INT UNSIGNED NOT NULL,
        FOREIGN KEY (category_id) REFERENCES categories (id),
        course_id INT UNSIGNED NOT NULL,
        FOREIGN KEY (course_id) REFERENCES course_details (id) ON DELETE CASCADE
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS course_general_settings;

-- +goose StatementEnd