-- +goose Up
CREATE TABLE IF NOT EXISTS quiz_attempt_sessions (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id INT UNSIGNED NOT NULL,
    student_id INT UNSIGNED NOT NULL,
    quiz_id INT UNSIGNED NOT NULL,
    attempt_number INT UNSIGNED NOT NULL DEFAULT 1,
    question_order JSON NOT NULL,
    started_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NULL DEFAULT NULL,
    submitted_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY idx_quiz_attempt_sessions_attempt (tenant_id, student_id, quiz_id, attempt_number),
    INDEX idx_quiz_attempt_sessions_quiz_id (quiz_id)
);

-- +goose Down
DROP TABLE IF EXISTS quiz_attempt_sessions;
