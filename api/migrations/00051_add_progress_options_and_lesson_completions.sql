-- +goose Up
ALTER TABLE course_certificate_settings
ADD COLUMN count_lessons BOOLEAN NOT NULL DEFAULT TRUE;

ALTER TABLE course_certificate_settings
ADD COLUMN count_quizzes BOOLEAN NOT NULL DEFAULT TRUE;

ALTER TABLE course_certificate_settings
ADD COLUMN count_assignments BOOLEAN NOT NULL DEFAULT TRUE;

-- +goose StatementBegin
CREATE TABLE
    student_lesson_completions (
        id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        tenant_id INT UNSIGNED NOT NULL,
        student_id INT UNSIGNED NOT NULL,
        course_id INT UNSIGNED NOT NULL,
        lesson_id INT UNSIGNED NOT NULL,
        completed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE,
        FOREIGN KEY (student_id) REFERENCES students (id) ON DELETE CASCADE,
        FOREIGN KEY (course_id) REFERENCES course_details (id) ON DELETE CASCADE,
        FOREIGN KEY (lesson_id) REFERENCES course_lessons (id) ON DELETE CASCADE,
        UNIQUE KEY uq_student_lesson_completions (tenant_id, student_id, lesson_id)
    );
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS student_lesson_completions;

ALTER TABLE course_certificate_settings DROP COLUMN count_assignments;

ALTER TABLE course_certificate_settings DROP COLUMN count_quizzes;

ALTER TABLE course_certificate_settings DROP COLUMN count_lessons;
