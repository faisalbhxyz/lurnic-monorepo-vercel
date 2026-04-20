-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    quiz_questions (
        id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        quiz_id INT UNSIGNED NOT NULL,
        CONSTRAINT fk_quiz FOREIGN KEY (quiz_id) REFERENCES course_quizzes (id) ON DELETE CASCADE,
        title VARCHAR(255) NOT NULL,
        details TEXT NULL DEFAULT NULL,
        media JSON NULL DEFAULT NULL,
        type ENUM ('multiple_choice', 'single_choice', 'true_false') NOT NULL DEFAULT 'single_choice',
        marks FLOAT DEFAULT 1,
        answer_required BOOLEAN NOT NULL DEFAULT FALSE,
        answer_explanation TEXT NULL DEFAULT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        INDEX idx_quiz_questions_quiz_id (quiz_id)
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE quiz_questions;

-- +goose StatementEnd