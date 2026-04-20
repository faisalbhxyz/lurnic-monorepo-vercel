-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    course_lessons (
        id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        description TEXT NULL,
        lesson_type ENUM ('video', 'live_session', 'audio', 'text') DEFAULT 'video',
        source_type ENUM (
            'youtube',
            'vimeo',
            'sound_cloud',
            'spotify',
            'custom_code',
            'upload'
        ) DEFAULT 'youtube',
        source JSON NULL,
        is_published BOOLEAN DEFAULT FALSE,
        is_public BOOLEAN DEFAULT FALSE,
        resources JSON NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        position INT DEFAULT 0,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        chapter_id INT UNSIGNED NOT NULL,
        FOREIGN KEY (chapter_id) REFERENCES course_chapters (id) ON DELETE CASCADE
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS course_lessons;

-- +goose StatementEnd