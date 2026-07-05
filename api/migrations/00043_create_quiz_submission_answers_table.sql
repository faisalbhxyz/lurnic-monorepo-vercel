-- +goose Up
CREATE TABLE IF NOT EXISTS quiz_submission_answers (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    submission_id INT UNSIGNED NOT NULL,
    question_id INT UNSIGNED NOT NULL,
    answer JSON NOT NULL,
    is_correct BOOLEAN NULL DEFAULT NULL,
    marks_awarded FLOAT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_quiz_submission_answers_submission_id (submission_id),
    INDEX idx_quiz_submission_answers_question_id (question_id)
);

-- +goose Down
DROP TABLE IF EXISTS quiz_submission_answers;
