-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    payment_methods (
        id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        title VARCHAR(100) NOT NULL,
        image TEXT NULL DEFAULT NULL,
        instruction TEXT NOT NULL,
        status BOOLEAN NOT NULL DEFAULT TRUE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        tenant_id INT UNSIGNED NOT NULL,
        FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payment_methods;

-- +goose StatementEnd