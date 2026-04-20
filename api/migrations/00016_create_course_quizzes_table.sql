-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    course_quizzes (
        id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        course_id INT UNSIGNED NOT NULL,
        CONSTRAINT fk_course FOREIGN KEY (course_id) REFERENCES course_details (id) ON DELETE CASCADE,
        chapter_id INT UNSIGNED NOT NULL,
        CONSTRAINT fk_chapter FOREIGN KEY (chapter_id) REFERENCES course_chapters (id) ON DELETE CASCADE,
        title VARCHAR(255) NOT NULL,
        instructions TEXT NOT NULL,
        is_published BOOLEAN NOT NULL DEFAULT FALSE,
        randomize_questions BOOLEAN NOT NULL DEFAULT FALSE,
        single_quiz_view BOOLEAN NOT NULL DEFAULT FALSE,
        time_limit INT NOT NULL DEFAULT 1,
        time_limit_option ENUM ('minutes', 'hours', 'days', 'weeks', 'months') NOT NULL DEFAULT 'weeks',
        total_visible_questions INT DEFAULT 0,
        reveal_answers BOOLEAN NOT NULL DEFAULT FALSE,
        enable_retry BOOLEAN NOT NULL DEFAULT FALSE,
        retry_attempts INT NOT NULL DEFAULT 1,
        minimum_pass_percentage FLOAT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        INDEX idx_course_quizzes_course_id (course_id),
        INDEX idx_course_quizzes_chapter_id (chapter_id)
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE course_quizzes;

-- +goose StatementEnd