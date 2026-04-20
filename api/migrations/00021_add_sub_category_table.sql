-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    sub_categories (
        id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        category_id INT UNSIGNED NOT NULL,
        FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE,
        name VARCHAR(255) NOT NULL,
        slug VARCHAR(255) NOT NULL,
        description TEXT NULL,
        thumbnail TEXT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        tenant_id INT UNSIGNED NOT NULL,
        FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sub_categories;

-- +goose StatementEnd