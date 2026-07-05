-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    student_lesson_video_progress (
        id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        tenant_id INT UNSIGNED NOT NULL,
        student_id INT UNSIGNED NOT NULL,
        course_id INT UNSIGNED NOT NULL,
        lesson_id INT UNSIGNED NOT NULL,
        max_position_seconds DOUBLE NOT NULL DEFAULT 0,
        duration_seconds DOUBLE NOT NULL DEFAULT 0,
        progress_percent DOUBLE NOT NULL DEFAULT 0,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE,
        FOREIGN KEY (student_id) REFERENCES students (id) ON DELETE CASCADE,
        FOREIGN KEY (course_id) REFERENCES course_details (id) ON DELETE CASCADE,
        FOREIGN KEY (lesson_id) REFERENCES course_lessons (id) ON DELETE CASCADE,
        UNIQUE KEY uq_student_lesson_video_progress (tenant_id, student_id, lesson_id)
    );
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS student_lesson_video_progress;
