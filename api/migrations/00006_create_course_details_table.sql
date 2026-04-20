-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    course_details (
        id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        summary TEXT NOT NULL,
        description TEXT NULL,
        visibility ENUM ('public', 'private', 'protected') NOT NULL DEFAULT 'private',
        is_scheduled BOOLEAN NOT NULL DEFAULT FALSE,
        schedule_date DATE NULL,
        schedule_time TIME NULL,
        featured_image VARCHAR(255) NULL,
        intro_video JSON NULL,
        pricing_model ENUM ('free', 'paid') NOT NULL DEFAULT 'free',
        regular_price DECIMAL(10, 2) NOT NULL DEFAULT 0,
        sale_price DECIMAL(10, 2) NOT NULL DEFAULT 0,
        show_comming_soon BOOLEAN NOT NULL DEFAULT FALSE,
        tags JSON NULL,
        overview JSON NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        author_id INT UNSIGNED NOT NULL,
        FOREIGN KEY (author_id) REFERENCES users (id) ON DELETE CASCADE,
        tenant_id INT UNSIGNED NOT NULL,
        FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS course_details;

-- +goose StatementEnd