-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    enrollments (
        id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        student_id INT UNSIGNED NOT NULL,
        FOREIGN KEY (student_id) REFERENCES students (id) ON DELETE CASCADE,
        course_id INT UNSIGNED NOT NULL,
        FOREIGN KEY (course_id) REFERENCES course_details (id) ON DELETE CASCADE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        tenant_id INT UNSIGNED NOT NULL,
        FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS enrollments;

-- +goose StatementEnd