-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    students (
        id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        user_id VARCHAR(255) UNIQUE NOT NULL,
        first_name VARCHAR(255) NOT NULL,
        last_name VARCHAR(255) NULL,
        phone VARCHAR(50) NULL,
        email VARCHAR(255) UNIQUE NOT NULL,
        password VARCHAR(255) NOT NULL,
        status BOOLEAN NOT NULL DEFAULT TRUE,
        otp_code VARCHAR(255) NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        tenant_id INT UNSIGNED NOT NULL,
        FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS students;

-- +goose StatementEnd