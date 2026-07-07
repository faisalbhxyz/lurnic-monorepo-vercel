-- +goose Up
CREATE TABLE IF NOT EXISTS assignment_attempt_sessions (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id INT UNSIGNED NOT NULL,
    student_id INT UNSIGNED NOT NULL,
    assignment_id INT UNSIGNED NOT NULL,
    started_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY idx_assignment_attempt_sessions_student (tenant_id, student_id, assignment_id),
    INDEX idx_assignment_attempt_sessions_assignment_id (assignment_id)
);

-- +goose Down
DROP TABLE IF EXISTS assignment_attempt_sessions;
