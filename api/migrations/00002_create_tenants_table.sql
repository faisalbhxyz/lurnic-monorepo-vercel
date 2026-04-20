-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS
    tenants (
        id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        app_key VARCHAR(255) UNIQUE NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    );

-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN IF NOT EXISTS tenant_id INT UNSIGNED NOT NULL;

-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE users
ADD CONSTRAINT fk_tenant FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP FOREIGN KEY fk_tenant;

-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN tenant_id;

-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS tenants;

-- +goose StatementEnd