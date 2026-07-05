-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    student_sessions (
        id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        session_id VARCHAR(255) NOT NULL,
        student_id INT UNSIGNED NOT NULL,
        tenant_id INT UNSIGNED NOT NULL,
        device_id VARCHAR(255) NOT NULL,
        device_name VARCHAR(255) NULL,
        user_agent VARCHAR(512) NULL,
        ip_address VARCHAR(45) NULL,
        last_seen_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        UNIQUE KEY uq_student_sessions_session_id (session_id),
        UNIQUE KEY uq_student_sessions_student_id (student_id),
        FOREIGN KEY (student_id) REFERENCES students (id) ON DELETE CASCADE,
        FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS student_sessions;
-- +goose StatementEnd
