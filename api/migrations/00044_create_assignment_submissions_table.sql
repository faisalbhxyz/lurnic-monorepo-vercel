-- +goose Up
-- +goose StatementBegin
CREATE TABLE assignment_submissions (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id INT UNSIGNED NOT NULL,
    course_id INT UNSIGNED NOT NULL,
    chapter_id INT UNSIGNED NOT NULL,
    assignment_id INT UNSIGNED NOT NULL,
    student_id INT UNSIGNED NOT NULL,
    response_text TEXT NULL DEFAULT NULL,
    score FLOAT NOT NULL DEFAULT 0,
    max_score FLOAT NOT NULL DEFAULT 0,
    percentage FLOAT NOT NULL DEFAULT 0,
    passed BOOLEAN NOT NULL DEFAULT FALSE,
    status ENUM ('submitted', 'graded', 'pending_review') NOT NULL DEFAULT 'pending_review',
    instructor_feedback TEXT NULL DEFAULT NULL,
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    graded_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_assignment_submissions_course_id (course_id),
    INDEX idx_assignment_submissions_assignment_id (assignment_id),
    INDEX idx_assignment_submissions_student_id (student_id),
    INDEX idx_assignment_submissions_tenant_id (tenant_id),
    UNIQUE KEY uniq_assignment_submissions_assignment_student (assignment_id, student_id)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE assignment_submissions;

-- +goose StatementEnd
