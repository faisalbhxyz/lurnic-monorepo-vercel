-- +goose Up
-- +goose StatementBegin
CREATE TABLE quiz_submissions (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id INT UNSIGNED NOT NULL,
    course_id INT UNSIGNED NOT NULL,
    chapter_id INT UNSIGNED NOT NULL,
    quiz_id INT UNSIGNED NOT NULL,
    student_id INT UNSIGNED NOT NULL,
    attempt_number INT UNSIGNED NOT NULL DEFAULT 1,
    score FLOAT NOT NULL DEFAULT 0,
    max_score FLOAT NOT NULL DEFAULT 0,
    percentage FLOAT NOT NULL DEFAULT 0,
    passed BOOLEAN NOT NULL DEFAULT FALSE,
    status ENUM ('submitted', 'graded', 'pending_review') NOT NULL DEFAULT 'submitted',
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    graded_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_quiz_submissions_course_id (course_id),
    INDEX idx_quiz_submissions_quiz_id (quiz_id),
    INDEX idx_quiz_submissions_student_id (student_id),
    INDEX idx_quiz_submissions_tenant_id (tenant_id)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE quiz_submissions;

-- +goose StatementEnd
