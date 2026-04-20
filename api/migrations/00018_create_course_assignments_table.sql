-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    course_assignments (
        id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        course_id INT UNSIGNED NOT NULL,
        CONSTRAINT fk_course_assignments_course FOREIGN KEY (course_id) REFERENCES course_details (id) ON DELETE CASCADE,
        chapter_id INT UNSIGNED NOT NULL,
        CONSTRAINT fk_course_assignments_chapter FOREIGN KEY (chapter_id) REFERENCES course_chapters (id) ON DELETE CASCADE,
        title VARCHAR(255) NOT NULL,
        instructions TEXT NOT NULL,
        attachments JSON NULL DEFAULT NULL,
        is_published BOOLEAN NOT NULL DEFAULT FALSE,
        time_limit INT NOT NULL DEFAULT 1,
        time_limit_option ENUM ('minutes', 'hours', 'days', 'weeks', 'months') NOT NULL DEFAULT 'weeks',
        file_upload_limit INT NOT NULL DEFAULT 1,
        total_marks FLOAT NOT NULL DEFAULT 1,
        minimum_pass_marks FLOAT NOT NULL DEFAULT 0,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        INDEX idx_course_assignments_course_id (course_id),
        INDEX idx_course_assignments_chapter_id (chapter_id)
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE course_assignments;

-- +goose StatementEnd