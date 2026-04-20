-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    banners (
        id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        title VARCHAR(100) NULL,
        url VARCHAR(255) NULL,
        image TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        tenant_id INT UNSIGNED NOT NULL,
        FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS banners;

-- +goose StatementEnd