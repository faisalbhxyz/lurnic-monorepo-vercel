-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    general_settings (
        id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        org_name VARCHAR(100) DEFAULT 'Lurnic',
        logo VARCHAR(255) NULL DEFAULT NULL,
        favicon VARCHAR(255) NULL DEFAULT NULL,
        student_prefix VARCHAR(10) DEFAULT 'S-',
        teacher_prefix VARCHAR(10) DEFAULT 'T-',
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        tenant_id INT UNSIGNED NOT NULL,
        CONSTRAINT fk_general_settings_tenant FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE,
        INDEX idx_general_settings_tenant_id (tenant_id)
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS general_settings;

-- +goose StatementEnd