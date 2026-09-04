-- +goose Up
-- +goose StatementBegin
CREATE TABLE student_class_profiles (
    student_id INT UNSIGNED NOT NULL PRIMARY KEY,
    class_level VARCHAR(32) NOT NULL,
    hsc_batch VARCHAR(32) NULL,
    department VARCHAR(32) NULL,
    preferred_class_slug VARCHAR(64) NULL,
    onboarding_completed TINYINT(1) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_scp_preferred_slug (preferred_class_slug),
    CONSTRAINT fk_scp_student FOREIGN KEY (student_id) REFERENCES students (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS student_class_profiles;
-- +goose StatementEnd
