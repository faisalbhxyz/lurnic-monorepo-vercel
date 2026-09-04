-- +goose Up
-- +goose StatementBegin
CREATE TABLE student_free_lessons (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id INT UNSIGNED NOT NULL,
    student_id INT UNSIGNED NOT NULL,
    lesson_id INT UNSIGNED NOT NULL,
    course_id INT UNSIGNED NOT NULL,
    added_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uq_student_free_lesson (tenant_id, student_id, lesson_id),
    INDEX idx_student_free_lessons_student (tenant_id, student_id, added_at),
    INDEX idx_student_free_lessons_lesson (lesson_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE,
    FOREIGN KEY (student_id) REFERENCES students (id) ON DELETE CASCADE,
    FOREIGN KEY (lesson_id) REFERENCES course_lessons (id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES course_details (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS student_free_lessons;
-- +goose StatementEnd
