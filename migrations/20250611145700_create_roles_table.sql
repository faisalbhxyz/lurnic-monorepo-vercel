-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    roles (
        id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        tenant_id INT UNSIGNED NULL,
        CONSTRAINT fk_roles_tenant FOREIGN KEY (tenant_id) REFERENCES tenants (id)
    );

-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN role_id INT UNSIGNED NOT NULL,
ADD CONSTRAINT fk_users_role FOREIGN KEY (role_id) REFERENCES roles (id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP FOREIGN KEY fk_users_role;

-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN role_id;

-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS roles;

-- +goose StatementEnd